package internal

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
	"github.com/hayabusa-cloud/hybs-server/application/model/games"
	"time"
)

func OnetimeTokenCreateMongodb(ctx hybs.Ctx) {
	// check parameters
	var serverID = ctx.FormString("server_id")
	if serverID == "" {
		ctx.StatusBadRequest()
		ctx.SysLogf("server_id is nil")
		return
	}
	if ctx.ConstStringValue("mgo_id") == "" {
		ctx.StatusInternalServerError("Bad ConstParams")
		return
	}
	ctx.SetCtxValue("HayabusaID", ctx.FormString("hayabusa_id"))
	// generate onetime token
	var newOnetimeToken = &games.OnetimeToken{
		HayabusaID:  ctx.FormString("hayabusa_id"),
		ServerID:    serverID,
		Token:       ctx.FormString("onetime_token"),
		Permission:  uint8(ctx.FormInt("permission") & 0xf),
		GeneratedAt: ctx.Now()}
	// insert into database
	var _, err = ctx.Mongo().Collection("onetime_token").InsertOne(ctx.Context(), newOnetimeToken)
	if err != nil {
		ctx.SysLogf("mongo insert onetime token %s failed:%s", newOnetimeToken, err)
		ctx.StatusInternalServerError("DB Error")
		return
	}
	// log
	ctx.GameLog("/account/onetime-token", map[string]interface{}{
		"category":      "account",
		"server_id":     serverID,
		"onetime_token": newOnetimeToken.Token,
		"permission":    newOnetimeToken.Permission,
		"generated_at":  newOnetimeToken.GeneratedAt,
	})
	ctx.SetResponseValue("onetimeToken", newOnetimeToken.ToMsg())
}

func OnetimeTokenCreateRedis(ctx hybs.Ctx) {
	// check parameters
	var serverID = ctx.FormString("server_id")
	if serverID == "" {
		ctx.StatusBadRequest()
		ctx.SysLogf("server_id is nil")
		return
	}
	if ctx.ConstStringValue("redis_id") == "" {
		ctx.StatusInternalServerError("Bad ConstParams")
		return
	}
	ctx.SetCtxValue("HayabusaID", ctx.FormString("hayabusa_id"))
	// generate onetime token
	var newOnetimeToken = &games.OnetimeToken{
		HayabusaID:  ctx.FormString("hayabusa_id"),
		ServerID:    serverID,
		Token:       ctx.FormString("onetime_token"),
		Permission:  uint8(ctx.FormInt("permission") & 0xf),
		GeneratedAt: ctx.Now()}
	// insert into redis
	if err := ctx.Redis().HSet(ctx.Context(), "OnetimeToken:"+newOnetimeToken.Token,
		"HayabusaID", newOnetimeToken.HayabusaID,
		"ServerID", newOnetimeToken.ServerID,
		"Permission", newOnetimeToken.Permission,
		"GeneratedAt", newOnetimeToken.GeneratedAt.Unix()).Err(); err != nil {
		ctx.SysLogf("redis set onetime token failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	ctx.Redis().Expire(ctx.Context(), "OnetimeToken:"+newOnetimeToken.Token,
		time.Duration(games.OnetimeTokenExpireDuration)*+time.Minute)
	// log
	ctx.GameLog("/account/onetime-token", map[string]interface{}{
		"category":      "account",
		"server_id":     serverID,
		"onetime_token": newOnetimeToken.Token,
		"permission":    newOnetimeToken.Permission,
		"generated_at":  newOnetimeToken.GeneratedAt,
	})
	ctx.SetResponseValue("onetimeToken", newOnetimeToken.ToMsg())
}
