package api

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/mat"

	statsanal "github.com/yodeman/analyses-api/stats-analyses"
)

// Response format for regression request
type regressionResp struct {
	Coeffs string `json:"regression_coefficients"`
	Tstats string `json:"t-test statistics"`
	Error  string `json:"error"`
}

// Request format for regression queries.
type regressionRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
}

func (server *Server) linearRegression(ctx *gin.Context) {
	var resp regressionResp
	var req regressionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		resp.Error = errResponse(fmt.Errorf(
			"Error parsing `username` in request body.\n%w", err))
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}

	authPayload, err := getPayload(ctx)
	if err != nil {
		resp.Error = errResponse(
			fmt.Errorf("Error getting authentication payload.\n%w", err))
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	if req.Username != authPayload.Username {
		resp.Error = errResponse(
			fmt.Errorf("request `username` and auth `username` don't match."))
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	userFile, err := server.querier.GetFile(ctx, authPayload.Username)
	if err != nil {
		resp.Error = errResponse(fmt.Errorf("Error fetching user's file\n%w", err))
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(userFile.Data)
	if err != nil {
		resp.Error = errResponse(fmt.Errorf("Error decoding user's file\n%w", err))
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	var data mat.Dense
	err = data.UnmarshalBinary(decoded)
	if err != nil {
		resp.Error = errResponse(fmt.Errorf("Error decoding user's file\n%w", err))
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	coeffs, tstats, err := statsanal.LinearRegression(&data)
	if err != nil {
		resp.Error = errResponse(fmt.Errorf("Error during regression analysis\n%w", err))
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Coeffs = coeffs
	resp.Tstats = tstats
	ctx.JSON(http.StatusOK, resp)
}
