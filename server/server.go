package server

import "github.com/bobykurniawan11/starter-go-prisma/config"

func Init() {
	config := config.GetConfig()
	r := NewRouter()
	r.Run(config.GetString("server.port"))
}
