package example

import "log/slog"

type MyStruct struct {
	structure *slog.Logger
	function  func()
}

func New() *MyStruct {
	return &MyStruct{}
}

func (ms *MyStruct) Print() {
	var log slog.Logger

	log.Debug("")

	ms.structure.Info("Hello :)")
}

func (ms *MyStruct) Call() {
	ms.function()
}
