package plugin_npm

import (
	"encoding/json"
	"fmt"
	"github.com/sinlov-go/go-common-lib/pkg/string_tools"
	"github.com/sinlov-go/go-common-lib/pkg/struct_kit"
	"github.com/sinlov-go/go-git-tools/git_info"
	"github.com/sinlov-go/unittest-kit/unittest_random_kit"
	"github.com/woodpecker-kit/woodpecker-npm/internal/pkgJson"
	"github.com/woodpecker-kit/woodpecker-tools/wd_flag"
	"github.com/woodpecker-kit/woodpecker-tools/wd_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_log"
	"github.com/woodpecker-kit/woodpecker-tools/wd_short_info"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// globalRegistry defines the default NPM registry.
const globalRegistry = "https://registry.npmjs.org/"

func (p *NpmPlugin) ShortInfo() wd_short_info.WoodpeckerInfoShort {
	if p.wdShortInfo == nil {
		info2Short := wd_short_info.ParseWoodpeckerInfo2Short(*p.woodpeckerInfo)
		p.wdShortInfo = &info2Short
	}
	return *p.wdShortInfo
}

// SetWoodpeckerInfo
// also change ShortInfo() return
func (p *NpmPlugin) SetWoodpeckerInfo(info wd_info.WoodpeckerInfo) {
	var newInfo wd_info.WoodpeckerInfo
	_ = struct_kit.DeepCopyByGob(&info, &newInfo)
	p.woodpeckerInfo = &newInfo
	info2Short := wd_short_info.ParseWoodpeckerInfo2Short(newInfo)
	p.wdShortInfo = &info2Short
}

func (p *NpmPlugin) GetWoodPeckerInfo() wd_info.WoodpeckerInfo {
	return *p.woodpeckerInfo
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

	if p.Settings.Registry != "" {
		_, errParseRegistry := url.Parse(p.Settings.Registry)
		if errParseRegistry != nil {
			return fmt.Errorf("parse registry, error by [ %s ] err: %v", p.Settings.Registry, errParseRegistry)
		}
	}

	if p.Settings.Tag != "" {
		if string_tools.StringInArr(p.Settings.Tag, tagForceNotSupport) {
			return fmt.Errorf("not support tag name [ %s ], tag name must not be: %v", p.Settings.Tag, tagForceNotSupport)
		}

		wd_log.Debugf("npm-auto-prerelease config: %v", p.Settings.TagAutoPrerelease)
		if p.Settings.TagAutoPrerelease {
			errChangeVersionByTagAuto := p.changePreReleaseVersionByTagAuto(p.Settings.Folder, p.Settings.Tag)
			if errChangeVersionByTagAuto != nil {
				return fmt.Errorf("check package version by semver, %v", errChangeVersionByTagAuto)
			}
		}

		wd_log.Debugf("npm-force-tag config: %v", p.Settings.TagForceEnable)
		if p.Settings.TagForceEnable { // check tag force enable
			errCheckSemver := p.checkPackageVersionBySemver()
			if errCheckSemver != nil {
				return fmt.Errorf("check package version by semver, %v", errCheckSemver)
			}
		}
	}

	return nil
}

const (
	preReleaseVersionCodeSize = 8
)

// changePreReleaseVersionByTagAuto
// format:
//
//	`<major>.<minor>.<patch>-<prerelease tag>.<prerelease build number>`
//
//	will try use CI `CI_COMMIT_SHA` first, if not get will use repo git head hash, if not get will use random code
func (p *NpmPlugin) changePreReleaseVersionByTagAuto(targetPkgFolder string, prereleaseTag string) error {
	wd_log.Debugf("changePreReleaseVersionByTagAuto start at folder: %s", targetPkgFolder)

	// Verify package.json file exists
	packagePath := filepath.Join(targetPkgFolder, "package.json")
	info, err := os.Stat(packagePath)

	if os.IsNotExist(err) {
		return fmt.Errorf("changePreReleaseVersion no package.json at %s: %w", packagePath, err)
	}
	if info.IsDir() {
		return fmt.Errorf("changePreReleaseVersion the package.json at %s is a directory", packagePath)
	}

	fileReadOld, err := os.ReadFile(packagePath)
	if err != nil {
		return fmt.Errorf("could not read package.json at %s: %w", packagePath, err)
	}

	npmPackageOld := npmPackage{}
	errUnmarshal := json.Unmarshal(fileReadOld, &npmPackageOld)
	if errUnmarshal != nil {
		return fmt.Errorf("could unmarshal at: %s, %w", packagePath, errUnmarshal)
	}

	var prereleaseCode string

	// use CI_COMMIT_SHA first
	if p.wdShortInfo.Commit.Sha != "" {
		if len(p.wdShortInfo.Commit.Sha) > preReleaseVersionCodeSize {
			prereleaseCode = p.wdShortInfo.Commit.Sha[:preReleaseVersionCodeSize]
		} else {
			prereleaseCode = p.wdShortInfo.Commit.Sha
		}
	} else {
		_, errIsPathGitManagementRoot := git_info.IsPathGitManagementRoot(p.Settings.RootPath)
		if errIsPathGitManagementRoot == nil {
			headInfo, errRepositoryHeadByPath := git_info.RepositoryHeadByPath(p.Settings.RootPath)
			if errRepositoryHeadByPath == nil {
				headHash := headInfo.Hash().String()
				if headHash != "" {
					wd_log.Debugf("current git head hash: %s", headHash)
					if len(headHash) > preReleaseVersionCodeSize {
						prereleaseCode = headHash[:preReleaseVersionCodeSize]
					} else {
						prereleaseCode = headHash
					}
				} else {
					wd_log.Warnf("get git repository head HASH err: %s", p.Settings.RootPath)
					prereleaseCode = unittest_random_kit.RandomStr(preReleaseVersionCodeSize)
					wd_log.Warnf("just use random code: %s", prereleaseCode)
				}
			}
		} else {
			wd_log.Warnf("run path not git repository root: %s", p.Settings.RootPath)
			prereleaseCode = unittest_random_kit.RandomStr(preReleaseVersionCodeSize)
			wd_log.Warnf("just use random code: %s", prereleaseCode)
		}
	}

	wd_log.Debugf("prereleaseCode code: %s", prereleaseCode)
	newVersion := fmt.Sprintf("%s-%s.%s", npmPackageOld.Version, prereleaseTag, prereleaseCode)

	errReplaceVersion := pkgJson.ReplaceJsonVersionByLine(packagePath, newVersion)
	if errReplaceVersion != nil {
		return fmt.Errorf("changePreReleaseVersion replace version err: %v", errReplaceVersion)
	}

	wd_log.Infof("npm-auto-prerelease version to: %s", newVersion)

	return nil
}

//func argCheckInArr(mark string, target string, checkArr []string) error {
//	if !(string_tools.StringInArr(target, checkArr)) {
//		return fmt.Errorf("not support %s now [ %s ], must in %v", mark, target, checkArr)
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
		return fmt.Errorf("verify the same registry in values do not match ci.yml [settings.npm-registry] : %s and package.json [publishConfig.registry] : %s", p.Settings.Registry, npm.Config.Registry)
	}

	p.npm = npm

	// Write the npmrc file
	if errWriteNpmrc := p.writeNpmrc(); errWriteNpmrc != nil {
		return fmt.Errorf("could not create npmrc: %w", errWriteNpmrc)
	}

	if p.Settings.DryRun {
		wd_log.Verbosef("dry run, skip some biz code, more info can open debug by flag [ %s ]", wd_flag.EnvKeyPluginDebug)
		return nil
	}

	// Attempt authentication
	if errAuthenticate := p.authenticate(); errAuthenticate != nil {
		return fmt.Errorf("could not authenticate: %w", errAuthenticate)
	}

	// Determine whether to publish
	publish, errPublish := p.shouldPublishPackage()

	if errPublish != nil {
		return fmt.Errorf("could not determine if package should be published: %w", errPublish)
	}

	if publish {
		wd_log.Info("Publishing package")

		if p.Settings.TagForceEnable && !p.Settings.DryRun {
			wd_log.Infof("unpublish package %s@%s", p.npm.Name, p.npm.Version)
			errUnpublish := runCommand(unpublishCommand(&p.Settings, p.npm.Name, p.npm.Version), p.Settings.Folder)
			if errUnpublish != nil {
				wd_log.Warnf("unpublish package fail: %v", errUnpublish)
			}
		}

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

// SetMockUserHome
// mock user home path for test
func (p *NpmPlugin) SetMockUserHome(userHome string) {
	p.mockUserHome = userHome
}
