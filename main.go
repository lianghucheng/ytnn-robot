package main

import (
	"log"
	"os"
	"os/signal"
	"ytnn-robot/robot"

	llog "github.com/name5566/leaf/log"
)

func init() {
	logger, err := llog.New("debug", "", log.Lshortfile|log.LstdFlags)
	if err != nil {
		panic(err)
	}
	llog.Export(logger)
}

func main() {
	robot.InitHall()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case sig := <-c:
		llog.Release("closing down (signal: %v)", sig)

		robot.DestroyGame()
		robot.DestroyHall()
	}
}
