package platform

import (
	"bytes"
	hybs "github.com/hayabusa-cloud/hayabusa"
	"github.com/hayabusa-cloud/hybs-server/application/common"
	"github.com/hayabusa-cloud/hybs-server/application/model/platform"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// authHayabusaID is middleware for basic auth
func authHayabusaID() hybs.Middleware {
	return func(h hybs.ServiceHandler) hybs.ServiceHandler {
		var basicAuthPrefix = []byte("Bearer ")
		return func(ctx hybs.Ctx) {
			for {
				// check auth header
				auth := ctx.RequestHeader().Peek("Authorization")
				if !bytes.HasPrefix(auth, basicAuthPrefix) {
					break
				}
				ctx.SysLogf("auth header checked:%s", auth)
				userBase := ctx.CtxValue("UserBase")
				if userBase == nil {
					break
				}
				ctx.SysLogf("user base information checked:%s", userBase.(*hybs.UserBase).UserID)
				// load data from mongodb
				var hayabusaID = &platform.HayabusaID{}
				var c = ctx.Mongo("Platform").Collection("hayabusa_id")
				if err := c.FindOne(ctx.Context(), bson.M{
					"user_id": userBase.(*hybs.UserBase).UserID,
				}).Decode(hayabusaID); err != nil {
					if err == mongo.ErrNoDocuments {
						ctx.StatusBadRequest()
					} else {
						ctx.StatusInternalServerError()
						ctx.SysLogf("mongo find hayabusa_id failed:%s", err)
					}
					return
				}
				// check permission level
				if hayabusaID.Permission < common.UserPermissionNormal {
					ctx.StatusForbidden("Permission Denied")
					return
				}
				// check ban
				if !hayabusaID.BanUntil.IsZero() && hayabusaID.BanUntil.After(ctx.Now()) {
					ctx.StatusForbidden("Account Banned")
					return
				}
				// set log fields
				ctx.SetCtxValue("LogFields", map[string]interface{}{
					"hayabusa_id":        hayabusaID.UserID,
					"user_created_day":   hayabusaID.CreatedAt.Format("2006-01-02"),
					"user_created_month": hayabusaID.CreatedAt.Format("2006-01"),
				})
				ctx.SetCtxValue("HayabusaID", hayabusaID)
				h(ctx)
				return
			}
			ctx.ResponseHeader().Set("WWW-Authenticate", "Bearer realm=Restricted")
			ctx.StatusUnauthorized()
		}
	}
}

func init() {
	hybs.RegisterMiddleware("AuthHayabusaID", authHayabusaID())
}
