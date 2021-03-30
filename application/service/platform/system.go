package platform

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
)

func systemTest(hybs.Ctx) {}

func init() {
	hybs.RegisterService("SystemTest", systemTest)
}
