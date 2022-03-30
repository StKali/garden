package api

import (
	db "github.com/stkali/garden/db/sqlc"
)
import "github.com/gin-gonic/gin"

type Server struct {
	store  db.Querier
	engine *gin.Engine
}

func NewServer(store db.Querier) *Server {
	server := &Server{
		store:  store,
		engine: gin.Default(),
	}
	registerRouts(server)
	return server
}

// registerRouts register routs to gin engine
func registerRouts(server *Server) {
	server.engine.POST("/user", server.CreateUser)
	server.engine.GET("/login", server.Login)
}
