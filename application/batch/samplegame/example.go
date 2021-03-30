package samplegame

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
)

func exampleTask1(ctx hybs.BatchCtx) {
}

func exampleTask2(ctx hybs.BatchCtx) {
}

func exampleTask3(ctx hybs.BatchCtx) {
}

func exampleTask4(ctx hybs.BatchCtx) {
}

func init() {
	hybs.RegisterBatchTask("SampleGameTask1", exampleTask1)
	hybs.RegisterBatchTask("SampleGameTask2", exampleTask2)
	hybs.RegisterBatchTask("SampleGameTask3", exampleTask3)
	hybs.RegisterBatchTask("SampleGameTask4", exampleTask4)
}
