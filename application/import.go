package application

import (
	_ "github.com/hayabusa-cloud/hybs-server/application/batch/samplegame"
	_ "github.com/hayabusa-cloud/hybs-server/application/middleware/admin"
	_ "github.com/hayabusa-cloud/hybs-server/application/middleware/games"
	_ "github.com/hayabusa-cloud/hybs-server/application/middleware/platform"
	_ "github.com/hayabusa-cloud/hybs-server/application/middleware/realtime"
	_ "github.com/hayabusa-cloud/hybs-server/application/service/admin"
	_ "github.com/hayabusa-cloud/hybs-server/application/service/central"
	_ "github.com/hayabusa-cloud/hybs-server/application/service/platform"
	_ "github.com/hayabusa-cloud/hybs-server/application/service/samplegame"
)
