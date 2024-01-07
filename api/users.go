package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"

	db "github.com/yodeman/analyses-api/dbase/sqlc"
	"github.com/yodeman/analyses-api/util"
)

// Response format for user queries.
// userResp hides sensitive information about user.
type userResp struct {
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}
type userResponse struct {
	User  userResp `json:"user"`
	Error string   `json:"error"`
}

// Request format for user queries.
type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
}

/*
createUser creates new user in the database using the data in the body
of the request. The endpoint expects a POST request with a json body
with the following key:

	`username`  - alphanumeric user's username
	`password`  - atleat 8 character user's password
	`email`     - user's email

The request returns response with the following http status codes:

200 - status OK:

	with response body:
	    {
	        "user": {
	            "username":"****",
	            "email": "*****",
	            "password_changed_at": "*****",
	            "created_at": "*****"
	        },
	        "error":""
	     }

400 - status Bad Request:

	Error parsing request body.
	with response body:
	    {
	        "user": {},
	        "error": "*****"
	    }

403 - status Forbidden:

	If username or email already exists.
	with response body:
	    {
	        "user": {},
	        "error": "*****"
	    }

501 - status Internal Server Error:

	with response body:
	    {
	        "user": {},
	        "error": "*****"
	    }
*/
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	var resp userResponse

	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Error = errResponse(fmt.Errorf("Error parsing request body.\n%w", err))
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		resp.Error = errResponse(
			fmt.Errorf("Error hashing user's password.\n%w", err))
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	user, err := server.querier.CreateUser(
		ctx,
		db.CreateUserParams{
			Username:       req.Username,
			HashedPassword: hashedPassword,
			Email:          req.Email,
		},
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" {
				resp.Error = errResponse(
					fmt.Errorf("User with username or email already exists.\n%w", err))
				ctx.JSON(http.StatusForbidden, resp)
				return
			}
		}

		resp.Error = errResponse(fmt.Errorf("Error while creating user.\n%w", err))
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.User = userResp{
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, resp)

	return
}

// Response format for user login
type loginResponse struct {
	AccessToken string   `json:"access_token"`
	User        userResp `json:"user"`
	Error       string   `json:"error"`
}

// Request format for login request
type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
}

/*
loginUser logs in registered user using the data in the request body.
The endpoint expects a GET request with a json body
with the following key:

	`username`  - alphanumeric user's username
	`password`  - atleat 8 character user's password

The request returns response with the following http status codes:

200 - status OK:

	with response body:
	    {
	        "access_token": "*****",
	        "user": {
	            "username":"****",
	            "email": "*****",
	            "password_changed_at": "*****",
	            "created_at": "*****"
	        },
	        "error":""
	     }

400 - status Bad Request:

	Error parsing request body.
	with response body:
	    {
	        "access_token": "",
	        "user": {},
	        "error": "*****"
	    }

404 - status Not Found:

	If user with username does not exist.
	with response body:
	    {
	        "access_token": "",
	        "user": {},
	        "error": "*****"
	    }

401 - status Unauthorized:

	If password does not match existing user's password.
	with response body:
	    {
	        "access_token": "",
	        "user": {},
	        "error": "*****"
	    }

501 - status Internal Server Error:

	with response body:
	    {
	        "user": {},
	        "error": "*****"
	    }
*/
func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	var resp loginResponse

	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Error = errResponse(fmt.Errorf("Error parsing request body.\n%w", err))
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	user, err := server.querier.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			resp.Error = errResponse(fmt.Errorf("User does not exist.\n%w", err))
			ctx.JSON(http.StatusNotFound, resp)
			return
		}

		resp.Error = errResponse(fmt.Errorf("Error fetching user.\n%w", err))
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		resp.Error = errResponse(fmt.Errorf("Incorrect password!\n%w", err))
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(
		req.Username,
		server.config.AccessTokenDuration)
	if err != nil {
		resp.Error = errResponse(fmt.Errorf("Error creating auth token.\n%w", err))
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.AccessToken = accessToken
	resp.User = userResp{
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, resp)

	return
}
