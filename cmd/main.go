package main

import (
	"flag"
	"fmt"
)

var (
	version   = "1.0.0"
	buildTime string
)

func main() {
	var (
		port int
		env  string
	)
	flag.IntVar(&port, "port", 4000, "API server port")
	flag.StringVar(&env, "env", "development", "Environment (development|staging|production)")

	flag.Parse()

	fmt.Println("wutuofu.com", port, env)
	fmt.Println(version, buildTime)
}
