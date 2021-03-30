package samplegame

import (
	"github.com/hayabusa-cloud/hybs-server/application/model/games"
)

type Player struct {
	games.PlayerBase `bson:",inline"`
	Nickname         string `bson:"nickname"`
	GoldNum          int    `bson:"gold_num"`
}

type MsgPlayer struct {
	HayabusaID string `json:"hayabusaId"`
	Nickname   string `json:"nickName"`
	GoldNum    int    `json:"goldNum"`
}

func (p *Player) ToMsg() *MsgPlayer {
	return &MsgPlayer{
		HayabusaID: p.HayabusaID,
		Nickname:   p.Nickname,
		GoldNum:    p.GoldNum,
	}
}
