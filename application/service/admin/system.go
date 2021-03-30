package admin

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
	"github.com/hayabusa-cloud/hybs-server/application/common"
	"strings"
)

// Admin API
type appListTableItem struct {
	AppName        string `json:"app_name"`
	AppDescription string `json:"app_description"`
	ServerNum      int    `json:"server_num"`
}

func adminAppsInfo(ctx hybs.Ctx) {
	data := make([]*appListTableItem, 0, len(common.Config.Apps))
	for _, app := range common.Config.Apps {
		data = append(data, &appListTableItem{
			AppName:        app.AppName,
			AppDescription: app.AppDescription,
			ServerNum:      len(app.Servers),
		})
	}
	ctx.SetResponseValue("data", data)
	ctx.SetResponseValue("size", len(data))
}

type serverListTableItem struct {
	ServerID          string `json:"server_id"`
	ServerDescription string `json:"server_description"`
}

func adminServersInfo(ctx hybs.Ctx) {
	data := make([]*serverListTableItem, 0)
	appName := ctx.FormString("app_name")
	if len(appName) > 0 {
		for _, app := range common.Config.Apps {
			if !strings.EqualFold(app.AppName, appName) {
				continue
			}
			for _, server := range app.Servers {
				data = append(data, &serverListTableItem{
					ServerID:          server.ID,
					ServerDescription: server.Description,
				})
			}
			break
		}
	}
	ctx.SetResponseValue("data", data)
	ctx.SetResponseValue("size", len(data))
}

type platformListTableItem struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Disabled    bool   `json:"disabled"`
	AppID       string `json:"app_id"`
	SDKUrl      string `json:"sdk_url"`
}

func adminPlatformsInfo(ctx hybs.Ctx) {
	data := make([]*platformListTableItem, 0)
	product := ctx.FormString("app_name")
	if len(product) > 0 {
		for _, app := range common.Config.Apps {
			if !strings.EqualFold(app.AppName, product) {
				continue
			}
			for _, platform := range app.Platforms {
				data = append(data, &platformListTableItem{
					Key:         platform.Key,
					Description: platform.Description,
					Disabled:    platform.Disabled,
					AppID:       platform.AppID,
				})
			}
			break
		}
	}
	ctx.SetResponseValue("data", data)
	ctx.SetResponseValue("size", len(data))
}

func init() {
	hybs.RegisterService("AdminAppsInfo", adminAppsInfo)
	hybs.RegisterService("AdminServersInfo", adminServersInfo)
	hybs.RegisterService("AdminPlatformsInfo", adminPlatformsInfo)
}
