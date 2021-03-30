package games

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
	"github.com/hayabusa-cloud/hybs-server/application/model/games"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
)

// authOnetimeTokenMongodb is middleware for basic auth and checking onetime token
// mongodb version
func authOnetimeTokenMongodb(h hybs.ServiceHandler) hybs.ServiceHandler {
	return func(ctx hybs.Ctx) {
		var appName = ctx.ConstStringValue("app_name")
		ctx.SetCtxValue("LogFields", map[string]interface{}{"app_name": appName})
		var reqOnetimeToken = ctx.FormString("onetime_token")
		if reqOnetimeToken == "" {
			ctx.SysLogf("onetime token cannot be null")
			ctx.StatusBadRequest()
			return
		}
		// in case of http method is not idempotent
		var newOnetimeToken, dbOnetimeToken = "", &games.OnetimeToken{}
		if !ctx.MethodIdempotent() {
			newOnetimeToken = games.NewOnetimeToken()
			var err = ctx.Mongo().Collection("onetime_token").FindOneAndUpdate(
				ctx.Context(),
				bson.M{"token": reqOnetimeToken},
				bson.M{"$set": bson.M{
					"token":        newOnetimeToken,
					"generated_at": ctx.Now().Raw(),
				}},
			).Decode(&dbOnetimeToken)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					ctx.SysLogf("mongo find onetime token %s failed:%s", reqOnetimeToken, err)
					ctx.StatusForbidden("Bad OnetimeToken")
				} else {
					ctx.SysLogf("mongo update onetime token %s failed:%s", reqOnetimeToken, err)
					ctx.StatusInternalServerError("DB Error")
				}
				return
			}
		} else {
			var err = ctx.Mongo().Collection("onetime_token").FindOne(
				ctx.Context(), bson.M{"token": reqOnetimeToken},
			).Decode(&dbOnetimeToken)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					ctx.SysLogf("mongo find onetime token %s failed:%s", reqOnetimeToken, err)
					ctx.StatusForbidden("Bad OnetimeToken")
				} else {
					ctx.SysLogf("mongo find onetime token %s failed:%s", reqOnetimeToken, err)
					ctx.StatusInternalServerError("DB Error")
				}
				return
			}
			newOnetimeToken = dbOnetimeToken.Token
		}
		// log
		ctx.SysLogf("authentication success:%s", dbOnetimeToken.HayabusaID)
		ctx.SysLogf("new onetime token=%s", newOnetimeToken)
		// set response
		dbOnetimeToken.Token = newOnetimeToken
		dbOnetimeToken.GeneratedAt = ctx.Now()
		ctx.SetCtxValue("HayabusaID", dbOnetimeToken.HayabusaID)
		ctx.SetCtxValue("UserBase", &hybs.UserBase{
			UserID:      dbOnetimeToken.HayabusaID,
			AccessToken: dbOnetimeToken.Token,
			Permission:  dbOnetimeToken.Permission,
		})
		ctx.SetResponseValue("onetimeToken", dbOnetimeToken.ToMsg())
		h(ctx)
		return
	}
}

// authOnetimeTokenRedis is middleware for basic auth and checking onetime token
// redis version
func authOnetimeTokenRedis(h hybs.ServiceHandler) hybs.ServiceHandler {
	return func(ctx hybs.Ctx) {
		var appName = ctx.ConstStringValue("app_name")
		ctx.SetCtxValue("LogFields", map[string]interface{}{"app_name": appName})
		var reqOnetimeToken = ctx.FormString("onetime_token")
		if reqOnetimeToken == "" {
			ctx.SysLogf("onetime token cannot be null")
			ctx.StatusBadRequest()
			return
		}
		// in case of http method is not idempotent
		var newOnetimeToken, dbOnetimeToken, cmdValues = "", &games.OnetimeToken{}, make([]interface{}, 0)
		if !ctx.MethodIdempotent() {
			newOnetimeToken = games.NewOnetimeToken()
			var scriptSHA1 = ctx.CtxString("/Redis/Scripts/OnetimeTokenGetSet")
			if scriptSHA1 == "" {
				ctx.SysLogf("redis eval script failed, SHA=%s not found", scriptSHA1)
				ctx.StatusInternalServerError()
				return
			}
			var cmd = ctx.Redis().EvalSha(
				ctx.Context(),
				scriptSHA1,
				// key: OnetimeToken:{OnetimeToken}
				[]string{"OnetimeToken:" + reqOnetimeToken, "OnetimeToken:" + newOnetimeToken},
				// val: GeneratedAt
				[]interface{}{-1, ctx.Now().Unix()}, // argv=generated at
			)
			if s, ok := cmd.Val().([]interface{}); ok {
				cmdValues = s
			}
			dbOnetimeToken.Token = newOnetimeToken
		} else {
			var scriptSHA1 = ctx.CtxString("/Redis/Scripts/OnetimeTokenGetOnly")
			if scriptSHA1 == "" {
				ctx.SysLogf("redis eval script failed, SHA=%s not found", scriptSHA1)
				ctx.StatusInternalServerError()
				return
			}
			var cmd = ctx.Redis().EvalSha(
				ctx.Context(),
				scriptSHA1,
				[]string{"OnetimeToken:" + reqOnetimeToken},
				[]interface{}{ctx.Now().Unix()})
			if s, ok := cmd.Val().([]interface{}); ok {
				cmdValues = s
			}
			dbOnetimeToken.Token = reqOnetimeToken
		}
		if len(cmdValues) != 4 {
			ctx.SysLogf("redis read onetime token failed")
			ctx.StatusInternalServerError()
			return
		}
		// unmarshal
		dbOnetimeToken.HayabusaID = cmdValues[0].(string)
		dbOnetimeToken.ServerID = cmdValues[1].(string)
		if permission, err := strconv.Atoi(cmdValues[2].(string)); err != nil {
			dbOnetimeToken.Permission = hybs.UserPermissionGuest
		} else {
			dbOnetimeToken.Permission = uint8(permission & 0xff)
		}
		if generatedAt, err := strconv.Atoi(cmdValues[3].(string)); err != nil {
			dbOnetimeToken.GeneratedAt = hybs.TimeNil()
		} else {
			dbOnetimeToken.GeneratedAt = hybs.TimeUnix(int64(generatedAt))
		}
		newOnetimeToken = dbOnetimeToken.Token

		// log
		ctx.SysLogf("authentication success:%s", dbOnetimeToken.HayabusaID)
		ctx.SysLogf("new onetime token=%s", newOnetimeToken)
		// set response
		ctx.SetCtxValue("HayabusaID", dbOnetimeToken.HayabusaID)
		ctx.SetCtxValue("UserBase", &hybs.UserBase{
			UserID:      dbOnetimeToken.HayabusaID,
			AccessToken: dbOnetimeToken.Token,
			Permission:  dbOnetimeToken.Permission,
		})
		ctx.SetResponseValue("onetimeToken", dbOnetimeToken.ToMsg())
		h(ctx)
		return
	}
}
func init() {
	// hybs.RegisterMiddleware("AuthOnetimeToken", authOnetimeTokenMongodb)
	hybs.RegisterMiddleware("AuthOnetimeToken", authOnetimeTokenRedis)
}
