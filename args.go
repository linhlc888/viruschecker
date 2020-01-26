package main

import (
	"flag"

	"github.com/olivere/env"
)

var (
	clamdHost = flag.String("clamd_host", env.String("127.0.0.1", "CLAMD_HOST"), "Clamd host")
	clamdPort = flag.Int("clamd_port", env.Int(3310, "CLAMD_PORT"), "Clamd port")
)

func init() {
	flag.Parse()
}
