package platform

import (
	"fmt"
	hybs "github.com/hayabusa-cloud/hayabusa"
	"github.com/hayabusa-cloud/hybs-server/application/common"
	"github.com/hayabusa-cloud/hybs-server/application/model/games"
	"github.com/hayabusa-cloud/hybs-server/application/model/platform"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"time"
)

func v1HayabusaIDCreate(ctx hybs.Ctx) {
	var newHayabusaID = platform.NewHayabusaID(ctx)
	newHayabusaID.Permission = common.UserPermissionNormal
	// insert new hayabusa id into mongodb
	var collection = ctx.Mongo("Platform").Collection("hayabusa_id")
	if _, err := collection.InsertOne(ctx.Context(), newHayabusaID); err != nil {
		ctx.SysLogf("mongo insert new platform account failed:%s", err)
		ctx.StatusInternalServerError("DB Error")
		return
	}
	// game log
	ctx.GameLog("/account/create", map[string]interface{}{
		"category":   "account",
		"user_id":    newHayabusaID.UserID,
		"permission": newHayabusaID.Permission,
		"created_at": newHayabusaID.CreatedAt,
	})
	// set response
	ctx.SetResponseValue("hayabusaID", newHayabusaID.UserID)
	ctx.SetResponseValue("accessToken", newHayabusaID.AccessToken)
}

func v1HayabusaIDGet(ctx hybs.Ctx) {
	var hayabusaID = ctx.CtxValue("HayabusaID")
	if hayabusaID == nil {
		ctx.StatusInternalServerError()
		return
	}
	ctx.SetResponseValue("hayabusaID", hayabusaID.(*platform.HayabusaID).ToMsg())
}

func v1OnetimeTokenCreate(ctx hybs.Ctx) {
	// check request params
	var val = ctx.CtxValue("HayabusaID")
	if val == nil {
		ctx.StatusInternalServerError()
		return
	}
	var hayabusaID = val.(*platform.HayabusaID)
	var appName = ctx.FormString("app_name")
	if appName == "" {
		ctx.SysLogf("app_name cannot be null")
		ctx.StatusBadRequest()
		return
	}
	var serverID = ctx.FormString("server_id")
	var tokenPermission = uint8(ctx.FormInt("permission") & 0xf)
	if tokenPermission > hayabusaID.Permission {
		ctx.SysLogf("token permission denied")
		ctx.StatusBadRequest()
		return
	}
	type msgEndpoint struct {
		Protocol string
		Address  string
		BasePath string
	}
	var onetimeTokenUrl, rootToken, msgGameServerEndpoint = "", "", msgEndpoint{}
	for _, app := range common.Config.Apps {
		if strings.EqualFold(appName, app.AppName) {
			for _, server := range app.Servers {
				if strings.EqualFold(serverID, server.ID) {
					msgGameServerEndpoint.Protocol = server.Scheme
					msgGameServerEndpoint.Address = server.Address
					msgGameServerEndpoint.BasePath = server.BasePath
					onetimeTokenUrl = server.OnetimeTokenUrl
					rootToken = server.RootToken
					break
				}
			}
			if onetimeTokenUrl != "" {
				break
			}
		}
	}
	if onetimeTokenUrl == "" {
		ctx.SysLogf("app_name=%s or server_id=%s not found", appName, serverID)
		ctx.StatusBadRequest()
		return
	}
	// check request times
	var generateFrequency = time.Duration(2 + len(common.Config.Apps))
	hayabusaID.Counter -= int(generateFrequency * ctx.Now().Sub(hayabusaID.CountedAt) / games.OnetimeTokenExpireDuration)
	if hayabusaID.Counter < 0 {
		hayabusaID.Counter = 0
	}
	if hayabusaID.Counter >= 15+len(common.Config.Apps) {
		ctx.Error(fasthttp.StatusTooManyRequests, "Request OnetimeToken Too Many Times")
		return
	}
	hayabusaID.Counter++
	hayabusaID.CountedAt = ctx.Now()
	// update hayabusa id counter value
	var _, err = ctx.Mongo("Platform").Collection("hayabusa_id").UpdateOne(
		ctx.Context(),
		bson.M{"user_id": hayabusaID.ID()},
		bson.M{"$set": bson.M{
			"counter":    hayabusaID.Counter,
			"counted_at": hayabusaID.CountedAt.Raw(),
		}})
	if err != nil {
		ctx.SysLogf("mongo update hayabusa_id failed:%s", err)
		ctx.StatusInternalServerError("DB Error")
		return
	}
	// request game server
	var onetimeToken = games.NewOnetimeToken()
	var req = fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(onetimeTokenUrl)
	req.Header.Set("Authorization", rootToken)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.PostArgs().Set("hayabusa_id", val.(*platform.HayabusaID).UserID)
	req.PostArgs().Set("server_id", serverID)
	req.PostArgs().Set("onetime_token", onetimeToken)
	req.PostArgs().Set("permission", fmt.Sprintf("%d", tokenPermission))
	var resp = fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if err = ctx.HttpClient().Do(req, resp); err != nil {
		ctx.SysLogf("post hayabusa_id onetime token to game server failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		ctx.SysLogf("post hayabusa_id onetime token to game server failed:%d", resp.StatusCode())
		ctx.StatusInternalServerError()
		return
	}
	var expireUntil = jsoniter.Get(resp.Body(), "onetimeToken", "expireUntil").ToInt64()
	// return response
	var msgOnetimeToken = &games.MsgOnetimeToken{
		Token:       onetimeToken,
		ExpireUntil: hybs.TimeUnix(expireUntil),
	}
	ctx.SetResponseValue("onetimeToken", msgOnetimeToken)
	ctx.SetResponseValue("endpoint", msgGameServerEndpoint)
}

func init() {
	hybs.RegisterService("PlatformCreateHayabusaIDV1", v1HayabusaIDCreate)
	hybs.RegisterService("PlatformGetHayabusaIDV1", v1HayabusaIDGet)
	hybs.RegisterService("PlatformCreateOnetimeTokenV1", v1OnetimeTokenCreate)
}
