package samplegame

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type PlayerGameToken struct {
	HayabusaID  string    `bson:"hayabusa_id"`
	GeneratedAt hybs.Time `bson:"generated_at"`
	Token       string    `bson:"token"`
}
type MsgPlayerGameToken struct {
	AppID      string    `json:"appId"`
	ValidUntil hybs.Time `json:"validUntil"`
	Token      string    `json:"token"`
}

func (token *PlayerGameToken) ToMsg(appID string) *MsgPlayerGameToken {
	return &MsgPlayerGameToken{
		AppID:      appID,
		ValidUntil: token.GeneratedAt.Add(time.Second * time.Duration(PlayerGameTokenPeriod) * 7 / 8),
		Token:      token.Token,
	}
}

// 秒数。mongoのSetExpireAfterSeconds()に使う
const PlayerGameTokenPeriod int32 = 60 * 60 * 24 * 8

var PlayerGameTokenIndex = mongo.IndexModel{
	Keys:    bson.D{{"hayabusa_id", -1}},
	Options: options.Index().SetName("idx_hayabusa_id").SetUnique(false),
}
var PlayerGameTokenExpiredIndex = mongo.IndexModel{
	Keys:    bson.D{{"generated_at", -1}},
	Options: options.Index().SetName("idx_expired").SetExpireAfterSeconds(PlayerGameTokenPeriod),
}
