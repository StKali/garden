package api

import (
	db "github.com/stkali/garden/db/sqlc"
	"github.com/stkali/garden/token"
	"github.com/stkali/garden/util"
	"github.com/gin-gonic/gin"
	
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
	server.engine.GET("/login", server.Login)
	server.engine.GET("/home", authMiddleware(server.maker), server.Home)
	server.engine.GET("", server.Home)
}

// Start http server on address
func (s *Server) Start(address string) {
	err := s.engine.Run(address)
	util.CheckError("cannot start HTTP server, err:", err)
}
