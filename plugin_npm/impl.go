package plugin_npm

import (
	"encoding/json"
	"fmt"
	"github.com/woodpecker-kit/woodpecker-tools/wd_log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// globalRegistry defines the default NPM registry.
const globalRegistry = "https://registry.npmjs.org/"

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
