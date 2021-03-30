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

func authorization(h hybs.RTHandler) hybs.RTHandler {
	return func(ctx hybs.RTCtx) {
		var cacheKey = fmt.Sprintf("/Realtime/Tokens/%s", ctx.Token())
		var val, ok = ctx.Cache().Get(cacheKey)
		if !ok {
			ctx.ErrorClientRequest(hybs.RTErrorCodeUnauthorized, "unauthorized")
			return
		}
		var token *AccessToken = nil
		if token, ok = val.(*AccessToken); !ok {
			ctx.ErrorServerInternal("access token wrong type")
			return
		}
		if !bytes.Equal(token.Token, ctx.Token()) {
			ctx.ErrorClientRequest(hybs.RTErrorCodeForbidden, "bad token")
			return
		}
		if ctx.Now().After(token.ValidUntil) {
			ctx.ErrorClientRequest(hybs.RTErrorCodeForbidden, "token expired")
			return
		}
		ctx.SetAuthorization(token)
		h(ctx)
	}
}

func init() {
	hybs.RegisterRTMiddleware("RTAuthorization", authorization)
}
