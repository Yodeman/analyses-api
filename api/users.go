package api

import (
    "database/sql"
    "errors"
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/lib/pq"

    "github.com/yodeman/analyses-api/util"
    db "github.com/yodeman/analyses-api/dbase/sqlc"
)

// Response format for user queries.
// userResp hides sensitive information about user.
type userResp struct {
    Username            string      `json:"username"`
    Email               string      `json:"email"`
    PasswordChangedAt   time.Time   `json:"password_changed_at"`
    CreatedAt           time.Time   `json:"created_at"`
}
type userResponse struct {
    User    userResp    `json:"user"`
    Error   string      `json:"error"`
}

// Request format for user queries.
type createUserRequest struct {
    Username    string  `json:"username" binding:"required,alphanum"`
    Password    string  `json:"password" binding:"required,min=8"`
    Email       string  `json:"email" binding:"required,email"`
}

// createUser creates new user in the database
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
                    fmt.Errorf("User with username and email already exists.\n%w", err))
                ctx.JSON(http.StatusForbidden, resp)
                return
            }
        }

        resp.Error = errResponse(fmt.Errorf("Error while creating user.\n%w", err))
        ctx.JSON(http.StatusInternalServerError, resp)
        return
    }

    resp.User = userResp{
        Username:           user.Username,
        Email:              user.Email,
        PasswordChangedAt:  user.PasswordChangedAt,
        CreatedAt:          user.CreatedAt,
    }
    ctx.JSON(http.StatusOK, resp)

    return
}

// Response format for user login
type loginResponse struct {
    AccessToken string      `json:"access_token"`
    User        userResp    `json:"user"`
    Error       string      `json:"error"`
}


// Request format for login request
type loginUserRequest struct {
    Username    string  `json:"username" binding:"required,alphanum"`
    Password    string  `json:"password" binding:"required,min=8"`
}

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
        Username:           user.Username,
        Email:              user.Email,
        PasswordChangedAt:  user.PasswordChangedAt,
        CreatedAt:          user.CreatedAt,
    }
    ctx.JSON(http.StatusOK, resp)

    return
}
