package gapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "simple-bank/db/sqlc"
	"simple-bank/pb"
	"simple-bank/token"
	"simple-bank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config      util.Config
	store       db.Store
	tokenMarker token.Maker
	router      *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	//tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:      config,
		store:       store,
		tokenMarker: tokenMaker,
	}

	return server, nil
}
