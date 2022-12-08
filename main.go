package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"time"

	"flag"
	"github.com/coreos/go-systemd/daemon"
)

func main() {
	typePtr := flag.String("systemd-type", "simple", "systemd type")
	flag.Parse()

	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		fmt.Printf("I am sleeping %d\n", i)
	}
	fmt.Println("I am up now!")
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Panicf("cannot listen: %s", err)
	}

	if *typePtr == "notify" {
		daemon.SdNotify(false, daemon.SdNotifyReady)
	}
	http.Serve(l, nil)
}
