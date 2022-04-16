package api

import (
	db "github.com/stkali/garden/db/sqlc"
	"github.com/stkali/garden/token"
)
import "github.com/gin-gonic/gin"

type Server struct {
	store  db.Querier
	maker  token.Maker
	engine *gin.Engine
}

func NewServer(store db.Querier, maker token.Maker) *Server {
	server := &Server{
		store:  store,
		maker:  maker,
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
