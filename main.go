package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"ytnn-robot/conf"
	"ytnn-robot/robot"

	llog "github.com/name5566/leaf/log"
)

func main() {
	logger, err := llog.New("debug", conf.GetCfgGameInfo().LogPath, log.Lshortfile|log.LstdFlags)
	if err != nil {
		panic(err)
	}
	llog.Export(logger)
	robot.RobotNumber = flag.Int("robotN", 100, "robot number")
	flag.Parse()
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
