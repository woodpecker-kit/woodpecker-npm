package plugin_npm

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sinlov-go/go-common-lib/pkg/string_tools"
	"github.com/woodpecker-kit/woodpecker-tools/wd_flag"
	"github.com/woodpecker-kit/woodpecker-tools/wd_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
)

type (
	// NpmPlugin plugin_npm all config
	NpmPlugin struct {
		Name           string
		Version        string
		WoodpeckerInfo *wd_info.WoodpeckerInfo
		onlyArgsCheck  bool
		Settings       Settings

		npm *npmPackage

		FuncPlugin FuncPlugin `json:"-"`
	}
)

type FuncPlugin interface {
	OnlyArgsCheck()

	Exec() error

	loadStepsTransfer() error
	checkArgs() error
	saveStepsTransfer() error
}

func (p *NpmPlugin) Exec() error {
	errLoadStepsTransfer := p.loadStepsTransfer()
	if errLoadStepsTransfer != nil {
		return errLoadStepsTransfer
	}

	errCheckArgs := p.checkArgs()
	if errCheckArgs != nil {
		return fmt.Errorf("check args err: %v", errCheckArgs)
	}
	if p.onlyArgsCheck {
		wd_log.Info("only check args, skip do doBiz")
		return nil
	}

	err := p.doBiz()
	if err != nil {
		return err
	}
	errSaveStepsTransfer := p.saveStepsTransfer()
	if errSaveStepsTransfer != nil {
		return errSaveStepsTransfer
	}

	return nil
}

func (p *NpmPlugin) OnlyArgsCheck() {
	p.onlyArgsCheck = true
}

func (p *NpmPlugin) loadStepsTransfer() error {
	return nil
}

func (p *NpmPlugin) checkArgs() error {

	errCheck := argCheckInArr("build status", p.WoodpeckerInfo.CurrentInfo.CurrentPipelineInfo.CiPipelineStatus, pluginBuildStateSupport)
	if errCheck != nil {
		return errCheck
	}

	if p.Settings.Token == "" {
		if p.Settings.Username == "" {
			return fmt.Errorf("missing username, please set %s", CliNameNpmUsername)
		}
		if p.Settings.Email == "" {
			return fmt.Errorf("missing email, please set %s", CliNameNpmEmail)
		}
		if p.Settings.Password == "" {
			return fmt.Errorf("missing password, please set %s", CliNameNpmPassword)
		}
	} else {
		wd_log.Info("Token credentials being used")
	}

	return nil
}

func argCheckInArr(mark string, target string, checkArr []string) error {
	if !(string_tools.StringInArr(target, checkArr)) {
		return fmt.Errorf("not support %s now [ %s ], must in %v", mark, target, checkArr)
	}
	return nil
}

//func checkEnvNotEmpty(keys []string) error {
//	for _, env := range keys {
//		if os.Getenv(env) == "" {
//			return fmt.Errorf("check env [ %s ] must set, now is empty", env)
//		}
//	}
//	return nil
//}

// doBiz
//
//	replace this code with your plugin_npm implementation
func (p *NpmPlugin) doBiz() error {

	if p.Settings.Folder == "" {
		p.Settings.Folder = p.Settings.RootPath
		wd_log.Debug("Just use root path as npm publish folder")
	}

	// Verify package.json file
	npm, err := readPackageFile(p.Settings.Folder)
	if err != nil {
		return fmt.Errorf("invalid package.json: %w", err)
	}

	// Verify the same registry is being used
	if p.Settings.Registry == "" {
		p.Settings.Registry = globalRegistry
	}

	if strings.Compare(p.Settings.Registry, npm.Config.Registry) != 0 {
		return fmt.Errorf("verify the same registry used settings registry values do not match .drone.yml: %s package.json: %s", p.Settings.Registry, npm.Config.Registry)
	}

	p.npm = npm

	// Write the npmrc file
	if errWriteNpmrc := p.writeNpmrc(); err != nil {
		return fmt.Errorf("could not create npmrc: %w", errWriteNpmrc)
	}

	// Attempt authentication
	if errAuthenticate := p.authenticate(); err != nil {
		return fmt.Errorf("could not authenticate: %w", errAuthenticate)
	}

	// Determine whether to publish
	publish, errPublish := p.shouldPublishPackage()

	if errPublish != nil {
		return fmt.Errorf("could not determine if package should be published: %w", errPublish)
	}

	if p.Settings.DryRun {
		wd_log.Verbosef("dry run, skip some biz code, more info can open debug by flag [ %s ]", wd_flag.EnvKeyPluginDebug)
		return nil
	}

	if publish {
		wd_log.Info("Publishing package")
		if err = runCommand(publishCommand(&p.Settings), p.Settings.Folder); err != nil {
			return fmt.Errorf("could not publish package: %w", err)
		}
	} else {
		wd_log.Info("Not publishing package")
	}

	return nil
}

func (p *NpmPlugin) saveStepsTransfer() error {
	// remove or change this code

	if p.Settings.StepsOutDisable {
		wd_log.Debugf("steps out disable by flag [ %v ], skip save steps transfer", p.Settings.StepsOutDisable)
		return nil
	}
	return nil
}

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

	// write npmrc file
	home := "/root"
	currentUser, err := user.Current()
	if err == nil {
		home = currentUser.HomeDir
	}
	npmrcPath := path.Join(home, ".npmrc")

	return os.WriteFile(npmrcPath, []byte(f(&p.Settings)), 0644)
}

// npmrcContentsUsernamePassword creates the contents from a username and
// password
func npmrcContentsUsernamePassword(config *Settings) string {
	// get the base64 encoded string
	authString := fmt.Sprintf("%s:%s", config.Username, config.Password)
	encoded := base64.StdEncoding.EncodeToString([]byte(authString))

	// create the file contents
	return fmt.Sprintf("_auth = %s\nemail = %s", encoded, config.Email)
}

// / Writes npmrc contents when using a token
func npmrcContentsToken(config *Settings) string {
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
	return fmt.Sprintf("%s:_authToken=%s", registryString, config.Token)
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

// shouldPublishPackage
// determines if the package should be published
func (p *NpmPlugin) shouldPublishPackage() (bool, error) {
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
