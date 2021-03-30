package master

import (
	"sync"

	hybs "github.com/hayabusa-cloud/hayabusa"
)

const (
	MGlobalVarInitMaxStamina         uint16 = 100
	MGlobalVarStaminaRecoverNum      uint16 = 101
	MGlobalVarStaminaRecoverInterval uint16 = 102
	MGlobalVarInitGoldNum            uint16 = 103
)

type GlobalVar struct {
	ID    uint16
	Value int
}

var (
	PlayerInitGoldNum = 0
	PlayerMaxGoldNum  = 999999999
)

func init() {
	// register master table and define insert method
	hybs.PrepareLoadMasterData("sample-game", func(data *GlobalVar, container *sync.Map) {
		container.Store(data.ID, data)
		switch data.ID {
		case MGlobalVarInitGoldNum:
			PlayerInitGoldNum = data.Value
		}
	})
}
