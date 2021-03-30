package games

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
	"github.com/hayabusa-cloud/hybs-server/application/model/games"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func checkPlayerStatus(h hybs.ServiceHandler) hybs.ServiceHandler {
	return func(ctx hybs.Ctx) {
		var val = ctx.User()
		if val == nil {
			ctx.StatusUnauthorized()
			return
		}
		// find player base data from db
		var hayabusaID = val.ID()
		var playerBase = &games.PlayerBase{}
		var result = ctx.Mongo().Collection("player").FindOne(
			ctx.Context(),
			bson.M{
				"hayabusa_id": hayabusaID,
			},
		)
		if err := result.Decode(playerBase); err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.SysLogf("player is nil, must sign up before use it")
				ctx.StatusBadRequest()
				return
			} else {
				ctx.SysLogf("mongo find player %s failed:%s", hayabusaID, err)
				ctx.StatusInternalServerError()
				return
			}
		}
		// check platform
		var platform = ctx.FormString("platform")
		if platform == "" {
			platform = "unknown"
		}
		// set log fields
		ctx.SetCtxValue("LogFields", map[string]interface{}{
			"app_name":           ctx.ConstStringValue("app_name"),
			"platform":           platform,
			"hayabusa_id":        playerBase.HayabusaID,
			"user_created_day":   playerBase.SignUpAt.Format("2006-01-02"),
			"user_created_month": playerBase.SignUpAt.Format("2006-01"),
		})
		// check ban
		if !playerBase.BanUntil.IsZero() && playerBase.BanUntil.After(ctx.Now()) {
			ctx.StatusForbidden("Account Banned")
			return
		}
		ctx.SysLogf("check player status success")
		h(ctx)
	}
}

func init() {
	hybs.RegisterMiddleware("CheckPlayerStatus", checkPlayerStatus)
}
