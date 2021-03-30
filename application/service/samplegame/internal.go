package samplegame

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
)

func __debug__(ctx hybs.Ctx) {
}

func init() {
	hybs.RegisterService("__debug__", __debug__)
}
