package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bobykurniawan11/starter-go-prisma/config"
	"github.com/bobykurniawan11/starter-go-prisma/db"
	"github.com/bobykurniawan11/starter-go-prisma/server"
)

func main() {

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}

	environment := flag.String("e", "development", "")

	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}

	flag.Parse()
	config.Init(*environment)

	server.Init()
}
