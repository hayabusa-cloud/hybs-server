package admin

import (
	"bytes"
	hybs "github.com/hayabusa-cloud/hayabusa"
	"github.com/hayabusa-cloud/hybs-server/application/common"
)

// authRoot is middleware for basic auth
func authRoot() hybs.Middleware {
	return func(h hybs.ServiceHandler) hybs.ServiceHandler {
		var basicAuthPrefix = []byte("Bearer ")
		return func(ctx hybs.Ctx) {
			for {
				auth := ctx.RequestHeader().Peek("Authorization")
				if !bytes.HasPrefix(auth, basicAuthPrefix) {
					break
				}
				// check root user
				rootToken := common.Config.RootToken
				if len(rootToken) > 0xf && string(auth) == rootToken {
					ctx.SetCtxValue("User", &hybs.UserBase{
						UserID:     "root",
						Permission: common.UserPermissionAdmin,
					})
					ctx.SysLogf("using root permission")
					h(ctx)
					return
				}
				userBase := ctx.CtxValue("User").(*hybs.UserBase)
				if userBase.Permission < common.UserPermissionAdmin {
					break
				}
				// set log fields
				ctx.SetCtxValue("LogFields", map[string]interface{}{
					"platform": "Admin",
				})
				// sign in
				ctx.SetCtxValue("User", &userBase)
				h(ctx)
				return
			}
			ctx.ResponseHeader().Set("WWW-Authenticate", "Bearer realm=Restricted")
			ctx.StatusForbidden()
		}
	}
}

func init() {
	hybs.RegisterMiddleware("AuthRoot", authRoot())
}
