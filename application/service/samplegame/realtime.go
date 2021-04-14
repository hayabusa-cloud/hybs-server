package samplegame

import (
	"encoding/base64"
	"fmt"
	"time"

	hybs "github.com/hayabusa-cloud/hayabusa"
	"github.com/hayabusa-cloud/hybs-server/application/middleware/realtime"
	"github.com/labstack/gommon/random"
)

const (
	rtEventCodeSampleGameSum = 0x8200
)

// a sample API: calculate the sum of a + b
func v1SampleGameRTTest(ctx hybs.RealtimeCtx) {
	var a, b int16
	ctx.ReadInt16(&a).ReadInt16(&b)
	var out = ctx.OutPacket()
	out.SetEventCode(rtEventCodeSampleGameSum)
	out.WriteInt16(a + b)
	ctx.Response(out) // response to requested user
	// ctx.Send(destination, out) // send message to specified user
	// ctx.BroadcastRoom(out)     // broadcast to all users in the same room
	// ctx.BroadcastServer(out)   // broadcast to all online users
	// can also get hayabusa id like this:
	// ctx.UserID()
}

// a sample middleware
func rtMiddlewareSampleGameQPS(h hybs.RealtimeHandler) hybs.RealtimeHandler {
	// here is concurrency safe
	var (
		requestCounter   = 0
		requestCountedAt = hybs.TimeNil()
	)
	return func(ctx hybs.RealtimeCtx) {
		// here is concurrency safe
		h(ctx)
		// statistics query per second
		if requestCountedAt.IsZero() {
			requestCountedAt = ctx.Now()
		}
		requestCounter++
		if ctx.Now().Unix() != requestCountedAt.Unix() {
			fmt.Println("query per second:", requestCounter)
			requestCounter = 0
		}
		requestCountedAt = ctx.Now()
	}
}

// generate temptation access token
func v1RealtimeAccessTokenPost(ctx hybs.Ctx) {
	var appID = ctx.FormString("app_id")
	if appID == "" {
		ctx.SysLogf("empty app_id")
		ctx.StatusBadRequest()
		return
	}
	// straightly generate onetime user account of playground env
	// implement customized accounts&authorization system in your own app
	var (
		onetimeUserID   = []byte(random.String(8, random.Alphanumeric))
		onetimePassword = []byte(random.String(15, random.Alphanumeric))
	)
	var src = []byte(fmt.Sprintf("%s#%s#%s", appID, onetimeUserID, onetimePassword))
	var token = []byte(base64.StdEncoding.EncodeToString(src))
	var accessToken = &realtime.AccessToken{
		AppID:      appID,
		HayabusaID: onetimeUserID,
		Mux:        "playground",
		Token:      token,
		ValidUntil: ctx.Now().Add(time.Hour * 24 * 3),
	}
	var cacheKey = fmt.Sprintf("/Realtime/Tokens/%s", accessToken.Token)
	ctx.Cache().Set(cacheKey, accessToken, time.Hour*24*3)
	ctx.SetResponseValue("token", accessToken.Token)
	ctx.SetResponseValue("expireUntil", accessToken.ValidUntil)
}

func init() {
	hybs.RegisterRealtimeHandler("RTSampleGameTestV1", v1SampleGameRTTest)
	hybs.RegisterRealtimeMiddleware("RTSampleGameQPS", rtMiddlewareSampleGameQPS)

	hybs.RegisterService("RealtimeCreateAccessTokenV1", v1RealtimeAccessTokenPost)
}
