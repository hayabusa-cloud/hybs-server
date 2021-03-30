package central

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SysRealtimeEndpoint struct {
	Network    string `bson:"network" json:"network" default:"udp"`
	Protocol   string `bson:"protocol" json:"protocol" default:"kcp"`
	Host       string `bson:"host" json:"host"`
	Port       uint16 `bson:"port" json:"port"`
	TurboLevel uint8  `bson:"turbo_level" json:"-"`
	Mtu        int    `bson:"mtu" json:"mtu"`
	SndWnd     int    `bson:"snd_wnd" json:"sndWnd"`
	RcvWnd     int    `bson:"rcv_wnd" json:"rcvWnd"`
	NoDelay    int    `bson:"no_delay" json:"noDelay"`
	Interval   int    `bson:"interval" json:"interval"`
	Resend     int    `bson:"resend" json:"resend"`
	Nc         int    `bson:"nc" json:"nc"`

	RootToken string `bson:"root_token" json:"-"`
	TokenUrl  string `bson:"token_url" json:"-"`
	Weight    int16  `bson:"weight" json:"-"`
	Count     int16  `bson:"count" json:"-"`
}

type SysRealtimeServer struct {
	AppID                string                 `bson:"app_id" json:"appId"`
	Owner                string                 `bson:"owner" json:"owner"`
	Enabled              bool                   `bson:"enabled" json:"-"`
	ValidUntil           *hybs.Time             `bson:"valid_until" json:"validUntil"`
	ExpireAt             *hybs.Time             `bson:"expire_at" json:"-"`
	Endpoints            []*SysRealtimeEndpoint `bson:"endpoints" json:"-"`
	CurrentEndpointIndex int                    `bson:"current_endpoint_index" json:"-"`
}

var SysRealtimeServerAppIDIndex = mongo.IndexModel{
	Keys:    bson.D{{"app_id", -1}, {"enabled", -1}},
	Options: options.Index().SetName("idx_app_id_enabled").SetUnique(true),
}
var SysRealtimeServerOwnerIndex = mongo.IndexModel{
	Keys:    bson.D{{"owner", -1}, {"enabled", -1}},
	Options: options.Index().SetName("idx_owner_enabled").SetUnique(false),
}
var SysRealtimeServerExpireIndex = mongo.IndexModel{
	Keys:    bson.D{{"expire_at", -1}},
	Options: options.Index().SetName("idx_expire_at").SetExpireAfterSeconds(1),
}
