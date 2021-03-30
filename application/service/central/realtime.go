package central

import (
	"fmt"
	"github.com/google/uuid"
	hybs "github.com/hayabusa-cloud/hayabusa"
	"github.com/hayabusa-cloud/hybs-server/application/common"
	"github.com/hayabusa-cloud/hybs-server/application/model/central"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

func v1RealtimeServerGet(ctx hybs.Ctx) {
	const cacheExpired = time.Minute
	// check input param
	var appID = ctx.RouteParams().ByName("app-id")
	if appID == "" {
		ctx.StatusNotFound()
		return
	}
	var realtimeServer = central.SysRealtimeServer{}
	// endpoint information does not require high timeliness, cache them
	var cacheKey = fmt.Sprintf("/Central/Realtime/AppID/%s", appID)
	for {
		// cache find first
		if val, ok := ctx.Cache().Get(cacheKey); ok {
			if val == nil {
				ctx.StatusNotFound()
				return
			}
			realtimeServer, ok = val.(central.SysRealtimeServer)
			if ok {
				break
			}
		}
		// mongo find
		var collection = ctx.Mongo("MongoCentral").Collection("sys_realtime_server")
		var err = collection.FindOne(
			ctx.Context(),
			bson.M{
				"app_id":  appID,
				"enabled": true,
			}).
			Decode(&realtimeServer)
		if err == mongo.ErrNoDocuments {
			ctx.Cache().Set(cacheKey, nil, time.Minute)
			ctx.StatusNotFound()
			return
		} else if err != nil {
			ctx.StatusInternalServerError()
			return
		}
		break
	}
	// check validity
	if realtimeServer.ValidUntil != nil && realtimeServer.ValidUntil.Before(ctx.Now()) {
		ctx.StatusNotFound()
		return
	}
	// select realtime server endpoint
	if realtimeServer.CurrentEndpointIndex >= len(realtimeServer.Endpoints) {
		realtimeServer.CurrentEndpointIndex = 0
	}
	var currentIndex = realtimeServer.CurrentEndpointIndex
	for {
		var endpoint = realtimeServer.Endpoints[currentIndex]
		if endpoint.Count < endpoint.Weight {
			break
		}
		endpoint.Count = 0
		currentIndex = (currentIndex + 1) % len(realtimeServer.Endpoints)
		if currentIndex == realtimeServer.CurrentEndpointIndex {
			break
		}
	}
	realtimeServer.CurrentEndpointIndex = currentIndex
	var selectedEndpoint = realtimeServer.Endpoints[currentIndex]
	selectedEndpoint.Count++
	// remote get access token from realtime server
	var req = fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(selectedEndpoint.TokenUrl)
	req.Header.Set("Authorization", selectedEndpoint.RootToken)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.PostArgs().Set("app_id", appID)
	var resp = fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if err := ctx.HttpClient().Do(req, resp); err != nil {
		ctx.SysLogf("remote get access token from realtime server failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		ctx.SysLogf("remote get access token from realtime server failed:%s",
			fasthttp.StatusMessage(resp.StatusCode()))
		ctx.StatusInternalServerError()
		return
	}
	var token = jsoniter.Get(resp.Body(), "token").ToString()
	var expireUntil = jsoniter.Get(resp.Body(), "expireUntil").ToInt64()
	// store cache data
	ctx.Cache().Set(cacheKey, realtimeServer, cacheExpired)
	// set response
	ctx.SetResponseValue("endpoint", selectedEndpoint)
	ctx.SetResponseValue("token", token)
	ctx.SetResponseValue("expireUntil", expireUntil)
}

// 3か月無料試用枠のAppIDを発行する（東京リージョン）
func v1RealtimeServerPostFreeTrial(ctx hybs.Ctx) {
	const cacheExpired = time.Minute
	var region = ctx.FormString("region")
	var appID = fmt.Sprintf("jp:playgroud:free-trial:%s", uuid.New().String())
	var validUntil = ctx.Now().AddDate(0, 3, 1)
	var expireAt = ctx.Now().AddDate(0, 3, 1)
	var endpoints = make([]*central.SysRealtimeEndpoint, 0)
	for _, server := range common.Config.RealtimeServerResources {
		if !strings.EqualFold(server.Region, region) {
			continue
		}
		var url = fmt.Sprintf("%s://%s:%d%s",
			server.TokenAPIScheme,
			server.TokenAPIHost,
			server.TokenAPIPort,
			server.TokenAPIPath)
		endpoints = append(endpoints, &central.SysRealtimeEndpoint{
			Network:   server.Network,
			Protocol:  server.Protocol,
			Host:      server.Host,
			Port:      server.Port,
			Mtu:       server.Mtu,
			SndWnd:    server.SndWnd,
			RcvWnd:    server.RcvWnd,
			NoDelay:   server.NoDelay,
			Interval:  server.Interval,
			Resend:    server.Resend,
			Nc:        server.Nc,
			RootToken: server.RootToken,
			TokenUrl:  url,
			Weight:    100,
			Count:     0,
		})
	}
	if len(endpoints) < 1 {
		ctx.SysLogf("region=%s not found", region)
		ctx.StatusBadRequest("region")
		return
	}
	var newApplication = &central.SysRealtimeServer{
		AppID:                appID,
		Owner:                "free-trial",
		Enabled:              true,
		ValidUntil:           &validUntil,
		ExpireAt:             &expireAt,
		Endpoints:            endpoints,
		CurrentEndpointIndex: 0,
	}
	var collection = ctx.Mongo("MongoCentral").Collection("sys_realtime_server")
	if _, err := collection.InsertOne(ctx.Context(), newApplication); err != nil {
		ctx.SysLogf("mongo insert realtime server information failed:%s", err)
		ctx.StatusInternalServerError()
		return
	}
	var cacheKey = fmt.Sprintf("/Central/Realtime/AppID/%s", appID)
	ctx.Cache().Set(cacheKey, newApplication, cacheExpired)
	ctx.SetResponseValue("realtimeApp", newApplication)
}

func init() {
	hybs.RegisterService("CentralGetRealtimeServerV1", v1RealtimeServerGet)
	hybs.RegisterService("CentralPostFreeTrialRealtimeServerV1", v1RealtimeServerPostFreeTrial)
}
