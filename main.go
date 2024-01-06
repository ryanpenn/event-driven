package main

import (
	"event-driven-system/user"
	"fmt"
	"time"
)

func main() {
	user.CreateUser()
	user.DeleteUser()

	fmt.Println()
	user.SamplePublish()

	time.Sleep(time.Second)
}
