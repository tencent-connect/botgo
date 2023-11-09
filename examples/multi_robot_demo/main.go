package main

import (
	"fmt"
	"multi_robot/local"
	"multi_robot/processor"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()

	local.InitProcessRobot("./conf/robot.yaml",
		processor.RegisterProcessor, processor.NewProcessor, []interface{}{
			InteractionEventHandler(),
			C2CMessageEventHandler(),
			GroupATMessageEventHandler(),
		}...)

	<-done
}
