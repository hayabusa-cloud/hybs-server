package master

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
	"sync"
)

type LevelExp struct {
	ID    int
	Level int
	Exp   int
}

func init() {
	// register master table and define insert method
	hybs.PrepareLoadMasterData("sample-game", func(data *LevelExp, container *sync.Map) {
		container.Store(data.Level, data)
	})
}
