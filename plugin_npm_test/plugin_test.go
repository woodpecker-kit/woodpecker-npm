package plugin_npm_test

import (
	"encoding/json"
	"github.com/woodpecker-kit/woodpecker-npm/plugin_npm"
	"github.com/woodpecker-kit/woodpecker-tools/wd_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_log"
	"github.com/woodpecker-kit/woodpecker-tools/wd_mock"
	"github.com/woodpecker-kit/woodpecker-tools/wd_short_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_steps_transfer"
	"testing"
)

func TestCheckArgsPlugin(t *testing.T) {
	t.Log("mock NpmPlugin")
	p := mockPluginWithStatus(t, wd_info.BuildStatusSuccess)

	// statusSuccess
	var statusSuccess plugin_npm.NpmPlugin
	deepCopyByPlugin(&p, &statusSuccess)
	statusSuccess.Settings.Username = "foo"
	statusSuccess.Settings.Password = "bar"
	statusSuccess.Settings.Email = "bar"

	// statusNotSupport
	var statusNotSupport plugin_npm.NpmPlugin
	deepCopyByPlugin(&p, &statusNotSupport)
	statusNotSupport.WoodpeckerInfo = wd_mock.NewWoodpeckerInfo(
		wd_mock.WithCurrentPipelineStatus("not_support"),
	)

	// noArgsUsername
	var noArgsUsername plugin_npm.NpmPlugin
	deepCopyByPlugin(&p, &noArgsUsername)
	noArgsUsername.Settings.Username = ""

	tests := []struct {
		name              string
		p                 plugin_npm.NpmPlugin
		isDryRun          bool
		workRoot          string
		wantArgFlagNotErr bool
	}{
		{
			name:              "statusSuccess",
			p:                 statusSuccess,
			wantArgFlagNotErr: true,
		},
		{
			name: "statusNotSupport",
			p:    statusNotSupport,
		},
		{
			name: "noArgsUsername",
			p:    noArgsUsername,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.p.OnlyArgsCheck()
			errPluginRun := tc.p.Exec()
			if tc.wantArgFlagNotErr {
				if errPluginRun != nil {
					wdShotInfo := wd_short_info.ParseWoodpeckerInfo2Short(*tc.p.WoodpeckerInfo)
					wd_log.VerboseJsonf(wdShotInfo, "print WoodpeckerInfoShort")
					wd_log.VerboseJsonf(tc.p.Settings, "print Settings")
					t.Fatalf("wantArgFlagNotErr %v\np.Exec() error:\n%v", tc.wantArgFlagNotErr, errPluginRun)
					return
				}
			} else {
				if errPluginRun == nil {
					t.Fatalf("test case [ %s ], wantArgFlagNotErr %v, but p.Exec() not error", tc.name, tc.wantArgFlagNotErr)
				}
				t.Logf("check args error: %v", errPluginRun)
			}
		})
	}
}

func TestPlugin(t *testing.T) {
	t.Log("do NpmPlugin")
	if envCheck(t) {
		return
	}
	if envMustArgsCheck(t) {
		return
	}
	t.Log("mock NpmPlugin")
	p := mockPluginWithStatus(t, wd_info.BuildStatusSuccess)
	//wd_log.VerboseJsonf(p, "print plugin_npm info")

	t.Log("mock plugin_npm config")

	// remove or change this code

	// statusSuccess
	var statusSuccess plugin_npm.NpmPlugin
	deepCopyByPlugin(&p, &statusSuccess)

	// statusFailure
	var statusFailure plugin_npm.NpmPlugin
	deepCopyByPlugin(&p, &statusFailure)
	statusFailure.WoodpeckerInfo = wd_mock.NewWoodpeckerInfo(
		wd_mock.WithCurrentPipelineStatus(wd_info.BuildStatusFailure),
	)

	// tagPipeline
	var tagPipeline plugin_npm.NpmPlugin
	deepCopyByPlugin(&p, &tagPipeline)
	tagPipeline.WoodpeckerInfo = wd_mock.NewWoodpeckerInfo(
		wd_mock.WithFastMockTag("v1.0.0", "new tag"),
	)

	// pullRequestPipeline
	var pullRequestPipeline plugin_npm.NpmPlugin
	deepCopyByPlugin(&p, &pullRequestPipeline)
	pullRequestPipeline.WoodpeckerInfo = wd_mock.NewWoodpeckerInfo(
		wd_mock.WithFastMockPullRequest("1", "new pr", "feature-support", "main", "main"),
	)

	tests := []struct {
		name            string
		p               plugin_npm.NpmPlugin
		isDryRun        bool
		workRoot        string
		ossTransferKey  string
		ossTransferData interface{}
		wantErr         bool
	}{
		{
			name: "statusSuccess",
			p:    statusSuccess,
		},
		{
			name:     "statusFailure",
			p:        statusFailure,
			isDryRun: true,
		},
		{
			name:     "tagPipeline",
			p:        tagPipeline,
			isDryRun: true,
		},
		{
			name:     "pullRequestPipeline",
			p:        pullRequestPipeline,
			isDryRun: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.p.Settings.DryRun = tc.isDryRun
			if tc.workRoot != "" {
				tc.p.Settings.RootPath = tc.workRoot
				errGenTransferData := generateTransferStepsOut(
					tc.p,
					tc.ossTransferKey,
					tc.ossTransferData,
				)
				if errGenTransferData != nil {
					t.Fatal(errGenTransferData)
				}
			}
			err := tc.p.Exec()
			if (err != nil) != tc.wantErr {
				t.Errorf("FeishuPlugin.Exec() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}

func mockPluginWithStatus(t *testing.T, status string) plugin_npm.NpmPlugin {
	p := plugin_npm.NpmPlugin{
		Name:    mockName,
		Version: mockVersion,
	}
	// use env:PLUGIN_DEBUG
	p.Settings.Debug = valEnvPluginDebug
	p.Settings.TimeoutSecond = envTimeoutSecond
	p.Settings.RootPath = testGoldenKit.GetTestDataFolderFullPath()
	p.Settings.StepsTransferPath = wd_steps_transfer.DefaultKitStepsFileName

	// mock woodpecker info
	//t.Log("mockPluginWithStatus")
	woodpeckerInfo := wd_mock.NewWoodpeckerInfo(
		wd_mock.WithCurrentPipelineStatus(status),
	)
	p.WoodpeckerInfo = woodpeckerInfo

	// mock all config at here

	return p
}

func deepCopyByPlugin(src, dst *plugin_npm.NpmPlugin) {
	if tmp, err := json.Marshal(&src); err != nil {
		return
	} else {
		err = json.Unmarshal(tmp, dst)
		return
	}
}

func generateTransferStepsOut(plugin plugin_npm.NpmPlugin, mark string, data interface{}) error {
	_, err := wd_steps_transfer.Out(plugin.Settings.RootPath, plugin.Settings.StepsTransferPath, *plugin.WoodpeckerInfo, mark, data)
	return err
}
