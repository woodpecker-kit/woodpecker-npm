package plugin_npm

import (
	"github.com/urfave/cli/v2"
	"github.com/woodpecker-kit/woodpecker-tools/wd_flag"
	"github.com/woodpecker-kit/woodpecker-tools/wd_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_log"
	"github.com/woodpecker-kit/woodpecker-tools/wd_short_info"
)

const (
	CliNameNpmRegistry = "settings.npm-registry"
	EnvNameNpmRegistry = "PLUGIN_NPM_REGISTRY"

	CliNameNpmUsername = "settings.npm-username"
	EnvNameNpmUsername = "PLUGIN_NPM_USERNAME"

	CliNameNpmPassword = "settings.npm-password"
	EnvNameNpmPassword = "PLUGIN_NPM_PASSWORD"

	CliNameNpmEmail = "settings.npm-email"
	EnvNameNpmEmail = "PLUGIN_NPM_EMAIL"

	CliNameNpmToken = "settings.npm-token"
	EnvNameNpmToken = "PLUGIN_NPM_TOKEN"

	CliNameNpmTag = "settings.npm-tag"
	EnvNameNpmTag = "PLUGIN_NPM_TAG"

	CliNameNpmFolder = "settings.npm-folder"
	EnvNameNpmFolder = "PLUGIN_NPM_FOLDER"

	CliNameNpmScopedAccess = "settings.npm-scoped-access"
	EnvNameNpmScopedAccess = "PLUGIN_NPM_SCOPED_ACCESS"

	CliNameNpmSkipVerifySSL = "settings.npm-skip-verify-ssl"
	EnvNameNpmSkipVerifySSL = "PLUGIN_NPM_SKIP_VERIFY_SSL"

	CliNameSkipWhoAmi = "settings.npm-skip-whoami"
	EnvNameSkipWhoAmi = "PLUGIN_SKIP_WHOAMI"

	CliNameFailOnVersionConflict = "settings.npm-fail-on-version-conflict"
	EnvNameFailOnVersionConflict = "PLUGIN_FAIL_ON_VERSION_CONFLICT"
)

// GlobalFlag
// Other modules also have flags
func GlobalFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    CliNameNpmRegistry,
			Usage:   "NPM registry to use when publishing packages. if empty will use https://registry.npmjs.org/",
			EnvVars: []string{EnvNameNpmRegistry},
		},
		&cli.StringFlag{
			Name:    CliNameNpmUsername,
			Usage:   "NPM username to use when publishing packages.",
			EnvVars: []string{EnvNameNpmUsername},
		},
		&cli.StringFlag{
			Name:    CliNameNpmPassword,
			Usage:   "NPM password to use when publishing packages.",
			EnvVars: []string{EnvNameNpmPassword},
		},
		&cli.StringFlag{
			Name:    CliNameNpmToken,
			Usage:   "NPM token to use when publishing packages. if token is set, username and password will be ignored.",
			EnvVars: []string{EnvNameNpmToken},
		},

		&cli.StringFlag{
			Name:    CliNameNpmEmail,
			Usage:   "NPM email to use when publishing packages.",
			EnvVars: []string{EnvNameNpmEmail},
		},

		&cli.StringFlag{
			Name:    CliNameNpmTag,
			Usage:   "NPM tag to use when publishing packages. this will cover package.json version field.",
			EnvVars: []string{EnvNameNpmTag},
		},
		&cli.StringFlag{
			Name:    CliNameNpmFolder,
			Usage:   "NPM folder to use when publishing packages which must containing package.json. default will use workspace",
			EnvVars: []string{EnvNameNpmFolder},
		},
		&cli.StringFlag{
			Name:    CliNameNpmScopedAccess,
			Usage:   "NPM scoped package access",
			EnvVars: []string{EnvNameNpmScopedAccess},
		},

		&cli.BoolFlag{
			Name:    CliNameNpmSkipVerifySSL,
			Usage:   "disables ssl verification when communicating with the NPM registry.",
			EnvVars: []string{EnvNameNpmSkipVerifySSL},
		},
		&cli.BoolFlag{
			Name:    CliNameSkipWhoAmi,
			Usage:   "Skip npm whoami check",
			EnvVars: []string{EnvNameSkipWhoAmi},
		},
		&cli.BoolFlag{
			Name:    CliNameFailOnVersionConflict,
			Usage:   "fail NPM publish if version already exists in NPM registry",
			EnvVars: []string{EnvNameFailOnVersionConflict},
		},
	}
}

func HideGlobalFlag() []cli.Flag {
	return []cli.Flag{}
}

func BindCliFlags(c *cli.Context,
	debug bool,
	cliName, cliVersion string,
	wdInfo *wd_info.WoodpeckerInfo,
	rootPath,
	stepsTransferPath string, stepsOutDisable bool,
) (*NpmPlugin, error) {

	config := Settings{
		Debug:             debug,
		TimeoutSecond:     c.Uint(wd_flag.NameCliPluginTimeoutSecond),
		StepsTransferPath: stepsTransferPath,
		StepsOutDisable:   stepsOutDisable,
		RootPath:          rootPath,

		Registry: c.String(CliNameNpmRegistry),
		Username: c.String(CliNameNpmUsername),
		Password: c.String(CliNameNpmPassword),
		Email:    c.String(CliNameNpmEmail),
		Token:    c.String(CliNameNpmToken),

		Folder:                c.String(CliNameNpmFolder),
		SkipVerifySSL:         c.Bool(CliNameNpmSkipVerifySSL),
		SkipWhoAmI:            c.Bool(CliNameSkipWhoAmi),
		FailOnVersionConflict: c.Bool(CliNameFailOnVersionConflict),
		Tag:                   c.String(CliNameNpmTag),
		ScopedAccess:          c.String(CliNameNpmScopedAccess),
	}

	// set default TimeoutSecond
	if config.TimeoutSecond == 0 {
		config.TimeoutSecond = 10
	}

	wd_log.Debugf("args %s: %v", wd_flag.NameCliPluginTimeoutSecond, config.TimeoutSecond)

	infoShort := wd_short_info.ParseWoodpeckerInfo2Short(*wdInfo)

	p := NpmPlugin{
		Name:           cliName,
		Version:        cliVersion,
		woodpeckerInfo: wdInfo,
		wdShortInfo:    &infoShort,
		Settings:       config,
	}

	return &p, nil
}
