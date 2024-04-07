package plugin_npm

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/woodpecker-kit/woodpecker-tools/wd_log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

const npmRcFileBase = ".npmrc"

// / writeNpmrc creates a .npmrc in the folder for authentication
func (p *NpmPlugin) writeNpmrc() error {
	var f func(settings *Settings) string
	if p.Settings.Token == "" {
		wd_log.DebugJsonf(p.Settings, "Specified credentials")
		f = npmrcContentsUsernamePassword
	} else {
		wd_log.Debug("Token credentials being used")
		f = npmrcContentsToken
	}

	npmrcPath := path.Join(p.Settings.Folder, npmRcFileBase)

	if p.Settings.NpmRcUserHomeEnable {
		// write npmrc file
		home := "/root"
		if p.mockUserHome == "" {
			currentUser, err := user.Current()
			if err == nil {
				home = currentUser.HomeDir
			}
		} else {
			home = p.mockUserHome
		}
		npmrcPath = path.Join(home, npmRcFileBase)
	}

	wd_log.Debugf("Writing npmrc file to %s", npmrcPath)
	return os.WriteFile(npmrcPath, []byte(f(&p.Settings)), 0644)
}

// npmrcContentsUsernamePassword creates the contents from a username and
// password
func npmrcContentsUsernamePassword(config *Settings) string {
	// get the base64 encoded string
	authString := fmt.Sprintf("%s:%s", config.Username, config.Password)
	encoded := base64.StdEncoding.EncodeToString([]byte(authString))

	// create the file contents
	if config.Registry == globalRegistry {
		return fmt.Sprintf("_auth = %s\nemail = %s", encoded, config.Email)
	}
	registryString := parseRegistryString(config)

	return fmt.Sprintf("%s:_auth = %s\nemail = %s", registryString, encoded, config.Email)
}

// / Writes npmrc contents when using a token
func npmrcContentsToken(config *Settings) string {
	registryString := parseRegistryString(config)
	return fmt.Sprintf("%s:_authToken=%s", registryString, config.Token)
}

func parseRegistryString(config *Settings) string {
	registry, _ := url.Parse(config.Registry)
	registry.Scheme = "" // Reset the scheme to empty. This makes it so we will get a protocol relative URL.
	host, port, _ := net.SplitHostPort(registry.Host)
	if port == "80" || port == "443" {
		registry.Host = host // Remove standard ports as they're not supported in authToken since NPM 7.
	}
	registryString := registry.String()

	if !strings.HasSuffix(registryString, "/") {
		registryString += "/"
	}
	return registryString
}

// authenticate
// try to auth with the NPM registry.
func (p *NpmPlugin) authenticate() error {
	var cmds []*exec.Cmd

	// Write the version command
	cmds = append(cmds, versionCommand())

	// write registry command
	if p.Settings.Registry != globalRegistry {
		cmds = append(cmds, registryCommand(p.Settings.Registry))
	}

	// Write skip verify command
	if p.Settings.SkipVerifySSL {
		cmds = append(cmds, skipVerifyCommand())
	}

	// Write whoami command to verify credentials
	if !p.Settings.SkipWhoAmI {
		cmds = append(cmds, whoamiCommand(p.Settings.Registry))
	}

	// Run commands
	err := runCommands(cmds, p.Settings.Folder)

	if err != nil {
		return err
	}

	return nil
}

// checkPackageVersionBySemver
func (p *NpmPlugin) checkPackageVersionBySemver() error {
	// Verify package.json file
	npm, err := readPackageFile(p.Settings.Folder)
	if err != nil {
		return fmt.Errorf("checkPackageVersionBySemver invalid package.json: %v", err)
	}
	targetVersion, errNpmVersion := semver.NewVersion(npm.Version)
	if errNpmVersion != nil {
		return fmt.Errorf("checkPackageVersionBySemver can not parse version: %s err: %v", npm.Version, errNpmVersion)
	}
	prereleaseInfo := targetVersion.Prerelease()
	if strings.Index(prereleaseInfo, p.Settings.Tag) != 0 {
		return fmt.Errorf("checkPackageVersionBySemver npm-tag [ %s ] must be the prefix of the prerelase version: [ %s ]", p.Settings.Tag, npm.Version)
	}
	return nil
}

// shouldPublishPackage
// determines if the package should be published
func (p *NpmPlugin) shouldPublishPackage() (bool, error) {

	if p.Settings.Tag != "" && p.Settings.TagForceEnable {
		wd_log.Debugf("skip check package version by semver, tag not empty and force enable")
		return true, nil
	}

	cmd := packageVersionsCommand(p.npm.Config.Registry, p.npm.Name)
	cmd.Dir = p.Settings.Folder

	trace(cmd)
	out, err := cmd.CombinedOutput()

	// see if there was an error
	// if there is an error its likely due to the package never being published
	if err == nil {
		// parse the json output
		var versions []string
		err = json.Unmarshal(out, &versions)

		if err != nil {
			wd_log.Debug("Could not parse into array of string. Likely single value")

			var version string
			errJson := json.Unmarshal(out, &version)

			if errJson != nil {
				return false, errJson
			}

			versions = append(versions, version)
		}

		for _, value := range versions {
			wd_log.Debugf("Found version of package: %s", value)

			if p.npm.Version == value {
				wd_log.Infof("Version found in the registry, as: %v", value)
				if p.Settings.FailOnVersionConflict {
					return false, fmt.Errorf("cannot publish package due to version conflict as version [ %v ]", value)
				}
				return false, nil
			}
		}

		wd_log.Info("Version not found in the registry")
	} else {
		wd_log.Info("Name was not found in the registry")
	}

	return true, nil
}

// / readPackageFile reads the package file at the given path.
func readPackageFile(folder string) (*npmPackage, error) {
	// Verify package.json file exists
	packagePath := filepath.Join(folder, "package.json")
	info, err := os.Stat(packagePath)

	if os.IsNotExist(err) {
		return nil, fmt.Errorf("no package.json at %s: %w", packagePath, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("the package.json at %s is a directory", packagePath)
	}

	// Read the file
	file, err := os.ReadFile(packagePath)
	if err != nil {
		return nil, fmt.Errorf("could not read package.json at %s: %w", packagePath, err)
	}

	// Unmarshal the json data
	npm := npmPackage{}
	err = json.Unmarshal(file, &npm)
	if err != nil {
		return nil, err
	}

	// Make sure values are present
	if npm.Name == "" {
		return nil, fmt.Errorf("no package name present")
	}
	if npm.Version == "" {
		return nil, fmt.Errorf("no package version present")
	}

	// Set the default registry
	if npm.Config.Registry == "" {
		npm.Config.Registry = globalRegistry
	}

	wd_log.DebugJsonf(npm, "npmPackage load as")
	return &npm, nil
}

// trace writes each command to standard error (preceded by a ‘$ ’) before it
// is executed. Used for debugging your build.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}

// runCommands executes the list of cmds in the given directory.
func runCommands(cmds []*exec.Cmd, dir string) error {
	for _, cmd := range cmds {
		err := runCommand(cmd, dir)

		if err != nil {
			return err
		}
	}

	return nil
}

func runCommand(cmd *exec.Cmd, dir string) error {
	wd_log.Debugf("runCommand start command\n%s\n", cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	trace(cmd)

	return cmd.Run()
}
