package main

import (
	"event-driven-system/user"
	"fmt"
)

func main() {
	user.CreateUser()
	fmt.Println()
	user.DeleteUser()
}
