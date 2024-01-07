package api

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gonum.org/v1/gonum/mat"

	db "github.com/yodeman/analyses-api/dbase/sqlc"
	"github.com/yodeman/analyses-api/util"
)

// Response format for file
// fileResp is used to hide information
type fileResp struct {
	ID        int64     `json:"id"`
	ChangedAt time.Time `json:"changed_at"`
}
type fileResponse struct {
	File  fileResp `json:"file"`
	Error string   `json:"error"`
}

/*
uploadFile uploads encoded user data to the database using the data in the body
of the request. If a user has a file in the database, the new file replaces the
old file in the database. The endpoint expects a POST request with a form-data
body with the following key:

	`username`   - alphanumeric user's username
	`file`       - a csv file.

The request returns response with the following http status codes:

200 - status OK:

	with response body:
	    {
	        "file": {
	            "id":"****",
	            "changed_at": "*****",
	        },
	        "error":""
	     }

400 - status Bad Request:

	Error parsing request body, and missing `username` or `file` key.
	with response body:
	    {
	        "file": {},
	        "error": "*****"
	    }

401 - status Unauthorized:

	If access token has expired or username in request body don't match username
	in token payload.
	with response body:
	    {
	        "file": {},
	        "error": "*****"
	    }

413 - status Request Entity Too Large:

	If file size exceed maximum limit.
	with response body:
	    {
	        "file": {},
	        "error": "*****"
	    }

501 - status Internal Server Error:

	with response body:
	    {
	        "file": {},
	        "error": "*****"
	    }
*/
func (server *Server) uploadFile(ctx *gin.Context) {
	var resp fileResponse

	file, err := ctx.FormFile("file")
	if err != nil {
		resp.Error = errResponse(
			fmt.Errorf("Error fetching data with key 'file'.\n%w", err))
		ctx.JSON(http.StatusBadRequest, resp)
		return
	}
	if file.Size > maxFileSize {
		resp.Error = errResponse(fmt.Errorf("File size > %d.\n", maxFileSize))
		ctx.JSON(http.StatusRequestEntityTooLarge, resp)
		return
	}

	username, ok := ctx.GetPostForm("username")
	if !ok {
		resp.Error = errResponse(
			fmt.Errorf("Error fetching username with key 'username'."))
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
	// check if the username matches with the token's username
	if username != authPayload.Username {
		resp.Error = errResponse(
			fmt.Errorf("request username and auth payload username mismatch!"))
		ctx.JSON(http.StatusUnauthorized, resp)
		return
	}

	reader, err := file.Open()
	if err != nil {
		resp.Error = errResponse(fmt.Errorf("Error opening uploaded file.\n%w", err))
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	rows, cols, data, err := util.ParseCSVToFloatSlice(reader)
	if err != nil {
		resp.Error = errResponse(fmt.Errorf("Error parsing uploaded file.\n%w", err))
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	bytes, err := mat.NewDense(rows, cols, data).MarshalBinary()
	if err != nil {
		resp.Error = errResponse(fmt.Errorf("Error encoding uploaded file.\n%w", err))
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	encoded := base64.StdEncoding.EncodeToString(bytes)

	userFile, err := server.querier.GetFile(ctx, username)
	if err != nil {
		// upload new file
		if errors.Is(err, sql.ErrNoRows) {
			userFile, err := server.querier.CreateFile(
				ctx,
				db.CreateFileParams{
					Username: username,
					Data:     encoded,
				},
			)
			if err != nil {
				resp.Error = errResponse(fmt.Errorf("Error uploading data.\n%w", err))
				ctx.JSON(http.StatusInternalServerError, resp)
				return
			}

			resp.File = fileResp{
				ID:        userFile.ID,
				ChangedAt: userFile.ChangedAt,
			}
			ctx.JSON(http.StatusOK, resp)
			return
		}
	}

	// file with user already exists, update the entry
	userFile, err = server.querier.UpdateFile(
		ctx,
		db.UpdateFileParams{
			Username: username,
			Data:     encoded,
		},
	)
	if err != nil {
		resp.Error = errResponse(fmt.Errorf("Error uploading data.\n%w", err))
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.File = fileResp{
		ID:        userFile.ID,
		ChangedAt: userFile.ChangedAt,
	}
	ctx.JSON(http.StatusOK, resp)
	return
}
