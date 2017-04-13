package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
)

var DockerLauncherSocket string = "/var/run/docker-launcher.sock"

func main() {
	logrus.Info("Container-Service Agent starting")
	fmt.Println("Start")

}
