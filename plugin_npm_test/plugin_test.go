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

	testCaseRootPath, errCreateTestCaseRootPath := testGoldenKit.GetOrCreateTestDataFullPath("check_args")
	if errCreateTestCaseRootPath != nil {
		t.Fatal(errCreateTestCaseRootPath)
	}

	// statusSuccess
	statusSuccessWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.FastWorkSpace(filepath.Join(testCaseRootPath, "statusSuccess")),
		wd_mock.FastCurrentStatus(wd_info.BuildStatusSuccess),
	)
	statusSuccessSettings := mockPluginSettings()
	statusSuccessSettings.Username = "foo"
	statusSuccessSettings.Password = "bar"
	statusSuccessSettings.Email = "baz"

	// tagNamLatestError
	tagNamLatestErrorWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.FastWorkSpace(filepath.Join(testCaseRootPath, "tagNamLatestError")),
		wd_mock.FastCurrentStatus(wd_info.BuildStatusSuccess),
	)
	tagNamLatestErrorSettings := mockPluginSettings()
	tagNamLatestErrorSettings.Username = "foo"
	tagNamLatestErrorSettings.Password = "bar"
	tagNamLatestErrorSettings.Email = "baz"
	tagNamLatestErrorSettings.Tag = "latest"

	// tagNamNextError
	tagNamNextErrorWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.FastWorkSpace(filepath.Join(testCaseRootPath, "tagNamNextError")),
		wd_mock.FastCurrentStatus(wd_info.BuildStatusSuccess),
	)
	tagNamNextErrorSettings := mockPluginSettings()
	tagNamNextErrorSettings.Username = "foo"
	tagNamNextErrorSettings.Password = "bar"
	tagNamNextErrorSettings.Email = "baz"
	tagNamNextErrorSettings.Tag = "latest"

	// tagForcePreReleaseError
	tagForcePreReleaseErrorWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.FastWorkSpace(filepath.Join(testCaseRootPath, "tagForcePreReleaseError")),
		wd_mock.FastCurrentStatus(wd_info.BuildStatusSuccess),
	)
	tagForcePreReleaseErrorSettings := mockPluginSettings()
	tagForcePreReleaseErrorSettings.Username = "foo"
	tagForcePreReleaseErrorSettings.Password = "bar"
	tagForcePreReleaseErrorSettings.Email = "baz"
	tagForcePreReleaseErrorSettings.TagForceEnable = true
	tagForcePreReleaseErrorSettings.Tag = "alpha"

	// tagForcePreReleaseRight
	tagForcePreReleaseRightWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.FastWorkSpace(filepath.Join(testCaseRootPath, "tagForcePreReleaseRight")),
		wd_mock.FastCurrentStatus(wd_info.BuildStatusSuccess),
	)
	tagForcePreReleaseRightSettings := mockPluginSettings()
	tagForcePreReleaseRightSettings.Username = "foo"
	tagForcePreReleaseRightSettings.Password = "bar"
	tagForcePreReleaseRightSettings.Email = "baz"
	tagForcePreReleaseRightSettings.TagForceEnable = true
	tagForcePreReleaseRightSettings.Tag = "alpha"

	// statusNotSupport
	statusNotSupportWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.FastWorkSpace(filepath.Join(testCaseRootPath, "statusNotSupport")),
		wd_mock.FastCurrentStatus("not_support"),
	)
	statusNotSupportSettings := mockPluginSettings()

	// noArgsUsername
	noArgsUsernameWoodpeckerInfo := *wd_mock.NewWoodpeckerInfo(
		wd_mock.FastWorkSpace(filepath.Join(testCaseRootPath, "noArgsUsername")),
		wd_mock.FastCurrentStatus(wd_info.BuildStatusSuccess),
	)
	noArgsUsernameSettings := mockPluginSettings()
	noArgsUsernameSettings.Username = ""

	tests := []struct {
		name           string
		woodpeckerInfo wd_info.WoodpeckerInfo
		settings       plugin_npm.Settings

		packageVersion string

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
			name:           "statusNotSupport",
			woodpeckerInfo: statusNotSupportWoodpeckerInfo,
			settings:       statusNotSupportSettings,
		},
		{
			name:           "noArgsUsername",
			woodpeckerInfo: noArgsUsernameWoodpeckerInfo,
			settings:       noArgsUsernameSettings,
		},
		{
			name:           "tagNamLatestError",
			woodpeckerInfo: tagNamLatestErrorWoodpeckerInfo,
			settings:       tagNamLatestErrorSettings,
		},
		{
			name:           "tagNamNextError",
			woodpeckerInfo: tagNamNextErrorWoodpeckerInfo,
			settings:       tagNamNextErrorSettings,
		},
		{
			name:              "tagForcePreReleaseRight",
			woodpeckerInfo:    tagForcePreReleaseRightWoodpeckerInfo,
			settings:          tagForcePreReleaseRightSettings,
			packageVersion:    "1.0.1-alpha.1",
			wantArgFlagNotErr: true,
		},
		{
			name:           "tagForcePreReleaseError",
			woodpeckerInfo: tagForcePreReleaseErrorWoodpeckerInfo,
			settings:       tagForcePreReleaseErrorSettings,
			packageVersion: "1.0.1",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.settings.Folder = tc.woodpeckerInfo.BasicInfo.CIWorkspace
			p := mockPluginWithSettings(t, tc.woodpeckerInfo, tc.settings)
			p.OnlyArgsCheck()

			if tc.packageVersion != "" {
				errMockPackageJsonFile := mockPackageJsonFile(p.Settings.RootPath, tc.name, tc.packageVersion, p.Settings.Registry)
				if errMockPackageJsonFile != nil {
					t.Fatal(errMockPackageJsonFile)
				}
			}

			errPluginRun := p.Exec()
			if tc.wantArgFlagNotErr {
				if errPluginRun != nil {
					wdShotInfo := wd_short_info.ParseWoodpeckerInfo2Short(p.GetWoodPeckerInfo())
					wd_log.VerboseJsonf(wdShotInfo, "print WoodpeckerInfoShort")
					wd_log.VerboseJsonf(p.Settings, "print Settings")
					t.Fatalf("wantArgFlagNotErr %v\np.Exec() error:\n%v", tc.wantArgFlagNotErr, errPluginRun)
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
				errMockPackageJsonFile := mockPackageJsonFile(p.Settings.RootPath, tc.name, "1.0.0", p.Settings.Registry)
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

func mockPackageJsonFile(root, pkgName, version string, registry string) error {
	pkgData := pkgJson.PkgJson{
		Name:    strings.ToLower(pkgName),
		Version: version,
		PublishConfig: pkgJson.NpmConfig{
			Registry: registry,
		},
	}
	pkgJsonPath := filepath.Join(root, "package.json")
	return unittest_file_kit.WriteFileAsJsonBeauty(pkgJsonPath, pkgData, true)
}
