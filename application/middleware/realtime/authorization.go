package realtime

import (
	"bytes"
	"fmt"

	hybs "github.com/hayabusa-cloud/hayabusa"
)

type AccessTokenInterface interface{}
type AccessToken struct {
	AppID      string
	HayabusaID []byte
	Mux        string
	Token      []byte
	ValidUntil hybs.Time
}

func authorization(h hybs.RealtimeHandler) hybs.RealtimeHandler {
	return func(ctx hybs.RealtimeCtx) {
		var cacheKey = fmt.Sprintf("/Realtime/Tokens/%s", ctx.Token())
		var val, ok = ctx.Cache().Get(cacheKey)
		if !ok {
			ctx.ErrorClientRequest(hybs.RealtimeErrorCodeUnauthorized, "unauthorized")
			return
		}
		var token *AccessToken = nil
		if token, ok = val.(*AccessToken); !ok {
			ctx.ErrorServerInternal("access token wrong type")
			return
		}
		if !bytes.Equal(token.Token, ctx.Token()) {
			ctx.ErrorClientRequest(hybs.RealtimeErrorCodeForbidden, "bad token")
			return
		}
		if ctx.Now().After(token.ValidUntil) {
			ctx.ErrorClientRequest(hybs.RealtimeErrorCodeForbidden, "token expired")
			return
		}
		ctx.SetAuthorization(token)
		h(ctx)
	}
}

func init() {
	hybs.RegisterRealtimeMiddleware("RTAuthorization", authorization)
}
