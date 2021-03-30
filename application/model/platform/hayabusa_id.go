package platform

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"regexp"
)

// HayabusaID defines structure of basic user data
type HayabusaID struct {
	// account data
	hybs.UserBase `bson:",inline" json:",inline"`
	CreatedAt     hybs.Time `json:"createdAt" bson:"created_at"`
	BanUntil      hybs.Time `json:"banUntil"  bson:"ban_until"`
	Counter       int       `json:"counter"   bson:"counter"`
	CountedAt     hybs.Time `json:"countedAt" bson:"counted_at"`
}

func (u *HayabusaID) ToMsg() *MsgHayabusaID {
	return &MsgHayabusaID{
		HayabusaID: u.UserID,
		Permission: u.Permission,
		CreatedAt:  u.CreatedAt,
	}
}

// NewHayabusaID returns new random HayabusaID with default config
func NewHayabusaID(ctx hybs.Ctx) *HayabusaID {
	var userBase = hybs.NewUser()
	return &HayabusaID{
		UserBase:  *userBase,
		CreatedAt: ctx.Now(),
		BanUntil:  hybs.TimeNil(),
		Counter:   0,
		CountedAt: ctx.Now(),
	}
}

// HayabusaIDIndex declares index of HayabusaID
var HayabusaIDIndex = mongo.IndexModel{
	Keys: bson.D{
		{"user_id", 1},
	},
	Options: options.Index().SetName("idx_hayabusa_id").SetUnique(true),
}

// MsgHayabusaID defines response structure of HayabusaID
type MsgHayabusaID struct {
	HayabusaID string    `json:"hayabusaId"`
	Permission uint8     `json:"permission"`
	CreatedAt  hybs.Time `json:"createdAt"`
}

// UserIDRegexp is default user id regexp rule
var UserIDRegexp = regexp.MustCompile("^[a-zA-Z0-9]{8,32}$")
