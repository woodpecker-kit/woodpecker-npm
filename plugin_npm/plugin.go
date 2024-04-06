package plugin_npm

import (
	"github.com/woodpecker-kit/woodpecker-tools/wd_info"
	"github.com/woodpecker-kit/woodpecker-tools/wd_short_info"
)

type (
	// NpmPlugin plugin_npm all config
	NpmPlugin struct {
		Name           string
		Version        string
		woodpeckerInfo *wd_info.WoodpeckerInfo
		wdShortInfo    *wd_short_info.WoodpeckerInfoShort
		onlyArgsCheck  bool
		Settings       Settings

		mockUserHome string

		npm *npmPackage

		FuncPlugin FuncPlugin `json:"-"`
	}
)

type FuncPlugin interface {
	ShortInfo() wd_short_info.WoodpeckerInfoShort

	SetWoodpeckerInfo(info wd_info.WoodpeckerInfo)
	GetWoodPeckerInfo() wd_info.WoodpeckerInfo

	OnlyArgsCheck()

	SetMockUserHome(userHome string)

	Exec() error

	loadStepsTransfer() error
	checkArgs() error
	saveStepsTransfer() error
}
