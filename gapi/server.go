package gapi

import (
	"fmt"
	db "simple-bank/db/sqlc"
	"simple-bank/pb"
	"simple-bank/token"
	"simple-bank/util"
	"simple-bank/worker"

	"github.com/gin-gonic/gin"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config      util.Config
	store       db.Store
	tokenMarker token.Maker
	router      *gin.Engine
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	//tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:      config,
		store:       store,
		tokenMarker: tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
