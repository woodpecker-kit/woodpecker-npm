package plugin_npm_test

import (
	"github.com/sinlov-go/unittest-kit/unittest_file_kit"
	"github.com/woodpecker-kit/woodpecker-npm/internal/pkgJson"
	"github.com/woodpecker-kit/woodpecker-npm/plugin_npm"
	"github.com/woodpecker-kit/woodpecker-tools/wd_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_log"
	"github.com/woodpecker-kit/woodpecker-tools/wd_mock"
	"github.com/woodpecker-kit/woodpecker-tools/wd_short_info"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckArgsPlugin(t *testing.T) {
	t.Log("mock NpmPlugin")

	// statusSuccess
	statusSuccessWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.FastCurrentStatus(wd_info.BuildStatusSuccess),
	)
	statusSuccessSettings := mockPluginSettings()
	statusSuccessSettings.Username = "foo"
	statusSuccessSettings.Password = "bar"
	statusSuccessSettings.Email = "baz"

	// registryError
	registryErrorWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.FastCurrentStatus(wd_info.BuildStatusSuccess),
	)
	registryErrorSettings := mockPluginSettings()
	registryErrorSettings.Registry = "some////foo.org"
	registryErrorSettings.Username = "foo"
	registryErrorSettings.Password = "bar"
	registryErrorSettings.Email = "baz"

	// statusNotSupport
	statusNotSupportWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.FastCurrentStatus("not_support"),
	)
	statusNotSupportSettings := mockPluginSettings()

	// noArgsUsername
	noArgsUsernameWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.FastCurrentStatus(wd_info.BuildStatusSuccess),
	)
	noArgsUsernameSettings := mockPluginSettings()
	noArgsUsernameSettings.Username = ""

	tests := []struct {
		name           string
		woodpeckerInfo wd_info.WoodpeckerInfo
		settings       plugin_npm.Settings

		isDryRun          bool
		wantArgFlagNotErr bool
	}{
		{
			name:              "statusSuccess",
			woodpeckerInfo:    statusSuccessWoodpeckerInfo,
			settings:          statusSuccessSettings,
			wantArgFlagNotErr: true,
		},
		{
			name:              "registryError",
			woodpeckerInfo:    registryErrorWoodpeckerInfo,
			settings:          registryErrorSettings,
			wantArgFlagNotErr: true,
		},
		{
			name:           "statusNotSupport",
			woodpeckerInfo: statusNotSupportWoodpeckerInfo,
			settings:       statusNotSupportSettings,
		},
		{
			name:           "noArgsUsername",
			woodpeckerInfo: noArgsUsernameWoodpeckerInfo,
			settings:       noArgsUsernameSettings,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := mockPluginWithSettings(t, tc.woodpeckerInfo, tc.settings)
			p.OnlyArgsCheck()
			errPluginRun := p.Exec()
			if tc.wantArgFlagNotErr {
				if errPluginRun != nil {
					wdShotInfo := wd_short_info.ParseWoodpeckerInfo2Short(p.GetWoodPeckerInfo())
					wd_log.VerboseJsonf(wdShotInfo, "print WoodpeckerInfoShort")
					wd_log.VerboseJsonf(p.Settings, "print Settings")
					t.Fatalf("wantArgFlagNotErr %v\np.Exec() error:\n%v", tc.wantArgFlagNotErr, errPluginRun)
					return
				}
				infoShot := p.ShortInfo()
				wd_log.VerboseJsonf(infoShot, "print WoodpeckerInfoShort")
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

	testCaseRootPath, errCreateTestCaseRootPath := testGoldenKit.GetOrCreateTestDataFullPath("plugin_npm")
	if errCreateTestCaseRootPath != nil {
		t.Fatal(errCreateTestCaseRootPath)
	}

	// statusSuccess
	statusSuccessWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.FastWorkSpace(filepath.Join(testCaseRootPath, "statusSuccess")),
		wd_mock.FastCurrentStatus(wd_info.BuildStatusSuccess),
	)
	statusSuccessSettings := mockPluginSettings()

	// statusFailure
	statusFailureWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.FastWorkSpace(filepath.Join(testCaseRootPath, "statusFailure")),
		wd_mock.FastCurrentStatus(wd_info.BuildStatusFailure),
	)
	statusFailureSettings := mockPluginSettings()

	// tagPipeline
	tagPipelineWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.FastWorkSpace(filepath.Join(testCaseRootPath, "tagPipeline")),
		wd_mock.FastTag("v1.0.0", "new tag"),
	)
	tagPipelineSettings := mockPluginSettings()

	// pullRequestPipeline
	pullRequestPipelineWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.FastWorkSpace(filepath.Join(testCaseRootPath, "pullRequestPipeline")),
		wd_mock.FastPullRequest("1", "new pr", "feature-support", "main", "main"),
	)
	pullRequestPipelineSettings := mockPluginSettings()

	tests := []struct {
		name           string
		woodpeckerInfo wd_info.WoodpeckerInfo
		settings       plugin_npm.Settings
		workRoot       string

		ossTransferKey  string
		ossTransferData interface{}

		isDryRun bool
		wantErr  bool
	}{
		{
			name:           "statusSuccess",
			woodpeckerInfo: statusSuccessWoodpeckerInfo,
			settings:       statusSuccessSettings,
			isDryRun:       true,
		},
		{
			name:           "statusFailure",
			woodpeckerInfo: statusFailureWoodpeckerInfo,
			settings:       statusFailureSettings,
			isDryRun:       true,
		},
		{
			name:           "tagPipeline",
			woodpeckerInfo: tagPipelineWoodpeckerInfo,
			settings:       tagPipelineSettings,
			isDryRun:       true,
		},
		{
			name:           "pullRequestPipeline",
			woodpeckerInfo: pullRequestPipelineWoodpeckerInfo,
			settings:       pullRequestPipelineSettings,
			isDryRun:       true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := mockPluginWithSettings(t, tc.woodpeckerInfo, tc.settings)
			p.Settings.DryRun = tc.isDryRun
			//p.SetMockUserHome(p.Settings.RootPath)
			if tc.ossTransferKey != "" {
				errGenTransferData := generateTransferStepsOut(
					p,
					tc.ossTransferKey,
					tc.ossTransferData,
				)
				if errGenTransferData != nil {
					t.Fatal(errGenTransferData)
				}
			}
			if p.Settings.Registry != "" {
				errMockPackageJsonFile := mockPackageJsonFile(p.Settings.RootPath, tc.name, p.Settings.Registry)
				if errMockPackageJsonFile != nil {
					t.Fatal(errMockPackageJsonFile)
				}
			}

			err := p.Exec()
			if (err != nil) != tc.wantErr {
				t.Errorf("FeishuPlugin.Exec() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}

func mockPackageJsonFile(root, pkgName string, registry string) error {
	pkgData := pkgJson.PkgJson{
		Name:    strings.ToLower(pkgName),
		Version: "1.0.0",
		PublishConfig: pkgJson.NpmConfig{
			Registry: registry,
		},
	}
	pkgJsonPath := filepath.Join(root, "package.json")
	return unittest_file_kit.WriteFileAsJsonBeauty(pkgJsonPath, pkgData, true)
}
