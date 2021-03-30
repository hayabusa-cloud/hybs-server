package samplegame

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
	"github.com/hayabusa-cloud/hybs-server/application/model/games"
	"github.com/hayabusa-cloud/hybs-server/application/model/samplegame"
	"github.com/hayabusa-cloud/hybs-server/application/model/samplegame/master"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func v1PlayerGet(ctx hybs.Ctx) {
	playerLoad(ctx)
	if ctx.HasError() {
		return
	}
	var player = ctx.CtxValue("Player").(*samplegame.Player)
	// refresh stamina point
	ctx.SetResponseValue("player", player.ToMsg())
}
func v1PlayerModify(ctx hybs.Ctx) {
	// get hayabusa id
	var playerID = ctx.User().ID()
	var mgoCollection = ctx.Mongo("MongoSampleGame").Collection("player")
	// default database id is defined by controller files, can be omitted as follow:
	// var mgoCollection = ctx.Mongo().Collection("player)

	// simple mongodb operation sample 1: select and unmarshal
	var player = &samplegame.Player{}
	var err = mgoCollection.FindOne(ctx.Context(),
		bson.M{
			"hayabusa_id": playerID,
		}).Decode(&player)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// requested player does not exist
			ctx.StatusBadRequest()
		} else {
			// server internal error
			ctx.SysLogf("mongo find player id=%s failed:%s", playerID, err)
			ctx.StatusInternalServerError()
		}
		return
	}
	// set response
	ctx.SetResponseValue("player", player.ToMsg())

	// simple mongodb operation sample 2: update
	var result = &mongo.UpdateResult{}
	result, err = mgoCollection.UpdateOne(ctx.Context(),
		bson.M{"hayabusa_id": playerID},
		bson.M{"$inc": bson.M{"gold_num": ctx.FormInt("gain_gold_num")}})
	if result.MatchedCount != 1 {
		ctx.StatusBadRequest()
		return
	}
	if err != nil {
		ctx.SysLogf("mongo update player failed:%s", err)
		return
	}

	// can write game log like this
	ctx.GameLog("/player/sample-action",
		map[string]interface{}{
			"player_id":    player.HayabusaID,
			"old_gold_num": player.GoldNum})
	// can write system log like this
	ctx.SysLogf("a sample system log item at %d", ctx.Now().Unix())
	// can response fields like this
	ctx.SetResponseValue("sampleIntField", 1000)
	ctx.SetResponseValue("sampleStringField", "hello hayabusa")
	ctx.SetResponseValue("sampleObjectField", struct {
		IntValue      int                    `json:"intValue"`
		BoolValue     bool                   `json:"boolValue"`
		StringValue   string                 `json:"stringValue"`
		DatetimeValue hybs.Time              `json:"datetimeValue"`
		SliceValue    []int                  `json:"sliceValue"`
		MapValue      map[string]interface{} `json:"mapValue"`
	}{
		IntValue:      500,
		BoolValue:     true,
		StringValue:   "hello hayabusa",
		DatetimeValue: ctx.Now(),
		SliceValue:    []int{10, 20, 30},
		MapValue:      map[string]interface{}{"a": 1, "b": "sample"},
	})
	// set status code like this
	ctx.StatusOk()
}

func init() {
	hybs.RegisterService("SampleGameGetPlayerV1", v1PlayerGet)
	hybs.RegisterService("SampleGameModifyPlayerV1", v1PlayerModify)
}

// ----- internal methods ----- //
func playerCreateInit(ctx hybs.Ctx, hayabusaID string) (player *samplegame.Player) {
	player = &samplegame.Player{
		PlayerBase: games.PlayerBase{
			HayabusaID: hayabusaID,
			SignUpAt:   ctx.Now(),
			BanUntil:   hybs.TimeNil(),
		},
		GoldNum: master.PlayerInitGoldNum,
	}
	return
}

// load player basic data
func playerLoad(ctx hybs.Ctx) {
	var hayabusaID = ctx.CtxString("HayabusaID")
	var player = &samplegame.Player{}
	if err := ctx.Mongo().Collection("player").FindOne(
		ctx.Context(),
		bson.M{"hayabusa_id": hayabusaID},
	).Decode(&player); err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.SysLogf("mongo find player failed:%s", err)
			ctx.StatusBadRequest()
		} else {
			ctx.SysLogf("mongo find player failed:%s", err)
			ctx.StatusInternalServerError("DB Error")
		}
		return
	}
	// set response
	ctx.SetCtxValue("Player", player)
	return
}

// sample business logic
func playerGoldModify(ctx hybs.Ctx, player *samplegame.Player, gold int, reason string) {
	ctx.SysLogf("gold modify %d", gold)
	if gold == 0 {
		return
	}
	if gold < 0 && player.GoldNum+gold < 0 {
		player.GoldNum = 0
	} else if gold > 0 && gold > master.PlayerMaxGoldNum-player.GoldNum {
		player.GoldNum = master.PlayerMaxGoldNum
	} else {
		player.GoldNum += gold
	}
	var _, err = ctx.Mongo().Collection("player").UpdateOne(
		ctx.Context(),
		bson.M{"hayabusa_id": player.HayabusaID},
		bson.M{"$set": bson.M{"gold_num": player.GoldNum}},
	)
	if err != nil && err == mongo.ErrNoDocuments {
		ctx.StatusBadRequest()
		return
	} else if err != nil {
		ctx.StatusInternalServerError()
		return
	}
}
