package games

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
	"github.com/labstack/gommon/random"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type OnetimeToken struct {
	HayabusaID  string    `bson:"hayabusa_id"`
	ServerID    string    `bson:"server_id"`
	Token       string    `bson:"token"`
	Permission  uint8     `bson:"permission"`
	GeneratedAt hybs.Time `bson:"generated_at"`
}

func (token *OnetimeToken) ToMsg() *MsgOnetimeToken {
	return &MsgOnetimeToken{
		Token:       token.Token,
		ExpireUntil: token.GeneratedAt.Add(time.Second * OnetimeTokenExpireDuration),
	}
}

type MsgOnetimeToken struct {
	Token       string    `json:"token"`
	ExpireUntil hybs.Time `json:"expireUntil"`
}

// NewOnetimeToken generates onetime token of user
func NewOnetimeToken() string {
	return random.String(32, random.Alphanumeric)
}

var OnetimeTokenHayabusaIndex = mongo.IndexModel{
	Keys:    bson.D{{"hayabusa_id", -1}},
	Options: options.Index().SetName("idx_hayabusa_id"),
}

var OnetimeTokenIndex = mongo.IndexModel{
	Keys:    bson.D{{"token", -1}},
	Options: options.Index().SetName("idx_onetime_token").SetUnique(true),
}

const OnetimeTokenExpireDuration = 60 * 60 * 6

var OnetimeTokenExpireIndex = mongo.IndexModel{
	Keys:    bson.D{{"generated_at", -1}},
	Options: options.Index().SetName("idx_expire").SetExpireAfterSeconds(OnetimeTokenExpireDuration + 60),
}
