package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ansible-semaphore/semaphore/cli/cmd"
)

func main() {
	log.Println("服务启动...")
	cmd.Execute()
}
