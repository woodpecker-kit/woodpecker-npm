package plugin_npm_test

import (
	"fmt"
	"github.com/sinlov-go/unittest-kit/env_kit"
	"github.com/sinlov-go/unittest-kit/unittest_file_kit"
	"github.com/woodpecker-kit/woodpecker-npm/plugin_npm"
	"github.com/woodpecker-kit/woodpecker-tools/wd_flag"
	"github.com/woodpecker-kit/woodpecker-tools/wd_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_log"
	"github.com/woodpecker-kit/woodpecker-tools/wd_steps_transfer"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

const (
	keyEnvDebug  = "CI_DEBUG"
	keyEnvCiNum  = "CI_NUMBER"
	keyEnvCiKey  = "CI_KEY"
	keyEnvCiKeys = "CI_KEYS"

	mockVersion = "v1.0.0"
	mockName    = "woodpecker-npm"
)

var (
	// testBaseFolderPath
	//  test base dir will auto get by package init()
	testBaseFolderPath = ""
	testGoldenKit      *unittest_file_kit.TestGoldenKit

	envTimeoutSecond uint

	// mustSetInCiEnvList
	//  for check set in CI env not empty
	mustSetInCiEnvList = []string{
		wd_flag.EnvKeyCiSystemPlatform,
		wd_flag.EnvKeyCiSystemVersion,
	}
	// mustSetArgsAsEnvList
	mustSetArgsAsEnvList = []string{
		plugin_npm.EnvNameNpmUsername,
	}

	valEnvPluginDebug = false

	valEnvRegistry    = ""
	valEnvNpmUsername = ""
	valEnvNpmPassword = ""
	valEnvNpmEmail    = ""
	valEnvNpmToken    = ""
)

func init() {
	testBaseFolderPath, _ = getCurrentFolderPath()
	wd_log.SetLogLineDeep(2)
	// if open wd_template please open this
	//wd_template.RegisterSettings(wd_template.DefaultHelpers)

	envTimeoutSecond = uint(env_kit.FetchOsEnvInt(wd_flag.EnvKeyPluginTimeoutSecond, 10))

	testGoldenKit = unittest_file_kit.NewTestGoldenKit(testBaseFolderPath)

	valEnvPluginDebug = env_kit.FetchOsEnvBool(wd_flag.EnvKeyPluginDebug, false)
	valEnvRegistry = env_kit.FetchOsEnvStr(plugin_npm.EnvNameNpmRegistry, "")
	valEnvNpmUsername = env_kit.FetchOsEnvStr(plugin_npm.EnvNameNpmUsername, "")
	valEnvNpmPassword = env_kit.FetchOsEnvStr(plugin_npm.EnvNameNpmPassword, "")
	valEnvNpmEmail = env_kit.FetchOsEnvStr(plugin_npm.EnvNameNpmEmail, "")
	valEnvNpmToken = env_kit.FetchOsEnvStr(plugin_npm.EnvNameNpmToken, "")
}

// test case basic tools start
// getCurrentFolderPath
//
//	can get run path this golang dir
func getCurrentFolderPath() (string, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("can not get current file info")
	}
	return filepath.Dir(file), nil
}

// test case basic tools end

func envCheck(t *testing.T) bool {

	if valEnvPluginDebug {
		wd_log.OpenDebug()
	}

	// most CI system will set env CI to true
	envCI := env_kit.FetchOsEnvStr("CI", "")
	if envCI == "" {
		t.Logf("not in CI system, skip envCheck")
		return false
	}
	t.Logf("check env for CI system")
	return env_kit.MustHasEnvSetByArray(t, mustSetInCiEnvList)
}

func envMustArgsCheck(t *testing.T) bool {
	for _, item := range mustSetArgsAsEnvList {
		if os.Getenv(item) == "" {
			t.Logf("plasee set env: %s, than run test\nfull need set env %v", item, mustSetArgsAsEnvList)
			return true
		}
	}
	return false
}

func generateTransferStepsOut(plugin plugin_npm.NpmPlugin, mark string, data interface{}) error {
	_, err := wd_steps_transfer.Out(plugin.Settings.RootPath, plugin.Settings.StepsTransferPath, plugin.GetWoodPeckerInfo(), mark, data)
	return err
}

func mockPluginSettings() plugin_npm.Settings {
	// all mock settings can set here
	settings := plugin_npm.Settings{
		// use env:PLUGIN_DEBUG
		Debug:             valEnvPluginDebug,
		TimeoutSecond:     envTimeoutSecond,
		RootPath:          testGoldenKit.GetTestDataFolderFullPath(),
		StepsTransferPath: wd_steps_transfer.DefaultKitStepsFileName,
	}

	settings.Registry = valEnvRegistry
	settings.Username = valEnvNpmUsername
	settings.Password = valEnvNpmPassword
	settings.Email = valEnvNpmEmail
	settings.Token = valEnvNpmToken

	return settings
}

func mockPluginWithSettings(t *testing.T, woodpeckerInfo wd_info.WoodpeckerInfo, settings plugin_npm.Settings) plugin_npm.NpmPlugin {
	p := plugin_npm.NpmPlugin{
		Name:    mockName,
		Version: mockVersion,
	}

	// mock woodpecker info
	//t.Log("mockPluginWithStatus")
	p.SetWoodpeckerInfo(woodpeckerInfo)

	if p.ShortInfo().Build.WorkSpace != "" {
		settings.RootPath = p.ShortInfo().Build.WorkSpace
	}

	p.Settings = settings
	return p
}
