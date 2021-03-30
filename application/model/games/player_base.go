package games

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PlayerBase struct {
	HayabusaID string    `bson:"hayabusa_id" json:"hayabusaId"` // アプリケーション層プライマリキー
	SignUpAt   hybs.Time `bson:"sign_up_at"  json:"signUpAt"`   // 登録日時
	BanUntil   hybs.Time `bson:"ban_until"   json:"banUntil"`   // 凍結期限
}

// PlayerIndex declares index of Player
var PlayerIndex = mongo.IndexModel{
	Keys:    bson.D{{"hayabusa_id", -1}},
	Options: options.Index().SetName("idx_hayabusa_id").SetUnique(true),
}
