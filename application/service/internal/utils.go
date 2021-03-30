package internal

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
	"github.com/hayabusa-cloud/hybs-server/application/common"
	"time"
)

func TimeParseDefault(value string) (ret hybs.Time, err error) {
	ret, err = hybs.TimeParse(common.TimeLayoutYMDHms, value)
	if err == nil {
		return
	}
	ret, err = hybs.TimeParse(common.TimeLayoutYMDHm, value)
	if err == nil {
		return
	}
	ret, err = hybs.TimeParse(time.RFC3339, value)
	if err == nil {
		return
	}
	return hybs.TimeParse(time.UnixDate, value)
}

// check platform
func CheckPlatform(pf string) (ok bool) {
	for _, app := range common.Config.Apps {
		for _, platform := range app.Platforms {
			if platform.Disabled {
				continue
			}
			if platform.Key != pf {
				continue
			}
			return true
		}
	}
	return false
}
