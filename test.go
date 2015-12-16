package main

import (
	"fmt"
	"time"
	"web/session"
	_ "web/session/memory"
)

func main() {
	session.NewManager("memory", "2", 3600)
	fmt.Println(time.Now().Unix())
}
