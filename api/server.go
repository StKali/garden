package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/stkali/garden/db/sqlc"
	"github.com/stkali/garden/token"
	"github.com/stkali/garden/util"
)

var setting = util.GetSetting()

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
	server.engine.POST("/login", server.Login)
	server.engine.GET("/home", authMiddleware(server.maker, server.store), server.Home)
	server.engine.GET("", authMiddleware(server.maker, server.store), server.Home)
}

// Start http server on address
func (s *Server) Start(address string) {
	err := s.engine.Run(address)
	util.CheckError("failed to start HTTP server", err)
}
