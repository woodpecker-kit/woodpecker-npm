package plugin_npm

import "os/exec"

// versionCommand gets the npm version
func versionCommand() *exec.Cmd {
	return exec.Command("npm", "--version")
}

// registryCommand sets the NPM registry.
func registryCommand(registry string) *exec.Cmd {
	return exec.Command("npm", "config", "set", "registry", registry)
}

// skipVerifyCommand disables ssl verification.
func skipVerifyCommand() *exec.Cmd {
	return exec.Command("npm", "config", "set", "strict-ssl", "false")
}

// whoamiCommand creates a command that gets the currently logged-in user.
func whoamiCommand(registry string) *exec.Cmd {
	if registry != "" {
		return exec.Command("npm", "whoami", "--registry", registry)
	} else {
		return exec.Command("npm", "whoami")
	}
}

// packageVersionsCommand gets the versions of the npm package.
func packageVersionsCommand(registry, name string) *exec.Cmd {
	return exec.Command("npm", "view", "--registry", registry, name, "versions", "--json")
}

// publishCommand runs the publish command
func publishCommand(settings *Settings) *exec.Cmd {
	commandArgs := []string{"publish"}

	if settings.Tag != "" {
		commandArgs = append(commandArgs, "--tag", settings.Tag)
	}

	if settings.ScopedAccess != "" {
		commandArgs = append(commandArgs, "--access", settings.ScopedAccess)
	}

	return exec.Command("npm", commandArgs...)
}
