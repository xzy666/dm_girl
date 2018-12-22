package main

import (
	"dm_girl/core"
	"os"
)

func main() {
	var player core.Player
	play := 1
	roomId := 1
	switch play {
	case core.BILIBILI:
		player.RecHandleFunc(core.BlHandler,roomId)
	//case core.DOUYU:
	//	player.RecHandleFunc(DyHandler)
	//case core.ONE:
	//	player.RecHandleFunc()
	//case core.TENCENT:
	//	player.RecHandleFunc()
	default:
		os.Exit(1000)
	}
	player.Connect()
}
