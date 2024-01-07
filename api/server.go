package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	db "github.com/yodeman/analyses-api/dbase/sqlc"
	"github.com/yodeman/analyses-api/token"
	"github.com/yodeman/analyses-api/util"
)

const maxFileSize = 10 << 20 // 10MB

type Server struct {
	config     util.Config
	querier    db.Querier
	router     *gin.Engine
	tokenMaker *token.PasetoMaker
}

func NewServer(config util.Config, querier db.Querier) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Error creating server.\n%w", err)
	}
	server := &Server{
		config:     config,
		querier:    querier,
		tokenMaker: tokenMaker,
	}

	router := gin.Default()
	router.MaxMultipartMemory = maxFileSize

	// request endpoints

	// user queries endpoints

	// register a user
	router.POST("/users/register", server.createUser)
	// login a user
	router.GET("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// users' files endpoints

	// upload file
	authRoutes.POST("/files/upload", server.uploadFile)

	// analyses endpoints

	// linear regression endpoint
	authRoutes.GET("/analyses/regression", server.linearRegression)

	server.router = router

	return server, nil
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

func errResponse(err error) string {
	return err.Error()
}

// getPayload gets authentication payload from request.
func getPayload(ctx *gin.Context) (*token.PasetoPayload, error) {
	payload, ok := ctx.MustGet(authorizationPayloadKey).(*token.PasetoPayload)
	if !ok {
		return nil, fmt.Errorf("Error getting payload from request.")
	}
	return payload, nil
}
