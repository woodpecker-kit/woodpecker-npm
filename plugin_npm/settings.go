package plugin_npm

import "github.com/woodpecker-kit/woodpecker-tools/wd_info"

const (
	// StepsTransferMarkDemoConfig
	// steps transfer key
	StepsTransferMarkDemoConfig = "demo_config"
)

var (
	// pluginBuildStateSupport
	pluginBuildStateSupport = []string{
		wd_info.BuildStatusCreated,
		wd_info.BuildStatusRunning,
		wd_info.BuildStatusSuccess,
		wd_info.BuildStatusFailure,
		wd_info.BuildStatusError,
		wd_info.BuildStatusKilled,
	}
)

type (
	// Settings plugin_npm private config
	Settings struct {
		Debug             bool
		TimeoutSecond     uint
		StepsTransferPath string
		StepsOutDisable   bool
		RootPath          string

		DryRun bool

		Registry string
		Username string
		Password string
		Email    string
		Token    string

		Folder                string
		SkipVerifySSL         bool
		SkipWhoAmI            bool
		FailOnVersionConflict bool
		Tag                   string
		ScopedAccess          string
	}

	npmPackage struct {
		Name    string    `json:"name"`
		Version string    `json:"version"`
		Config  npmConfig `json:"publishConfig"`
	}

	npmConfig struct {
		Registry string `json:"registry"`
	}
)
