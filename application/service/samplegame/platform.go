package samplegame

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
	"github.com/hayabusa-cloud/hybs-server/application/service/internal"
	"go.mongodb.org/mongo-driver/mongo"
)

// mongodb implements version
func v1OnetimeTokensPostMongodb(ctx hybs.Ctx) {
	// call common service
	internal.OnetimeTokenCreateMongodb(ctx)
	if ctx.HasError() {
		return
	}
	createNewPlayerIfNotExists(ctx)
}

// redis implements version
func v1OnetimeTokensPostRedis(ctx hybs.Ctx) {
	// call common service
	internal.OnetimeTokenCreateRedis(ctx)
	if ctx.HasError() {
		return
	}
	createNewPlayerIfNotExists(ctx)
}

func init() {
	// hybs.RegisterService("SampleGameCreateOnetimeTokenV1", v1OnetimeTokensPostMongodb)
	hybs.RegisterService("SampleGameCreateOnetimeTokenV1", v1OnetimeTokensPostRedis)
}

func createNewPlayerIfNotExists(ctx hybs.Ctx) {
	// init create player
	var hayabusaID = ctx.CtxValue("HayabusaID").(string)
	// 基本情報データの作成
	var player = playerCreateInit(ctx, hayabusaID)
	var collection = ctx.Mongo().Collection("player")
	if _, err := collection.InsertOne(ctx.Context(), player); err != nil {
		switch err.(type) {
		case mongo.WriteException:
			return
		default:
			ctx.StatusInternalServerError()
			ctx.SysLogf("mongo insert player failed:%s", err)
		}
		return
	}
}
