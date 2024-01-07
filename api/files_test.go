package api

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
	"gonum.org/v1/gonum/mat"

	mockdb "github.com/yodeman/analyses-api/dbase/mock"
	db "github.com/yodeman/analyses-api/dbase/sqlc"
	"github.com/yodeman/analyses-api/token"
	"github.com/yodeman/analyses-api/util"
)

func TestUploadFile(t *testing.T) {
	user, _ := randomUser(t)
	sampleCSV := util.RandomCSV(30, 10) // 30 rows and 10 cols csv file
	rows, cols, data, err := util.ParseCSVToFloatSlice(strings.NewReader(sampleCSV))
	require.NoError(t, err)
	matrix := mat.NewDense(rows, cols, data)
	byteData, err := matrix.MarshalBinary()
	require.NoError(t, err)
	encoded := base64.StdEncoding.EncodeToString(byteData)

	createFileParams := db.CreateFileParams{
		Username: user.Username,
		Data:     encoded,
	}

	uploadResp := db.File{
		ID:       util.RandomInt(1, 1000),
		Username: user.Username,
		Data:     encoded,
	}

	testCases := []struct {
		name       string
		params     map[string]string
		buildStubs func(querier *mockdb.MockQuerier)
		setupAuth  func(
			t *testing.T, request *http.Request, tokenMaker *token.PasetoMaker,
		)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "NEW FILE",
			params: map[string]string{
				"username":    user.Username,
				"usernameKey": "username",
				"fileKey":     "file",
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetFile(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.File{}, sql.ErrNoRows)

				querier.EXPECT().
					CreateFile(gomock.Any(), gomock.Eq(createFileParams)).
					Times(1).
					Return(uploadResp, nil)
			},
			setupAuth: func(
				t *testing.T, request *http.Request, tokenMaker *token.PasetoMaker,
			) {
				addAuthorization(
					t, request, tokenMaker, authorizationTypeToken, user.Username,
					time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchFile(t, recorder.Body, uploadResp)
			},
		},
		{
			name: "UPDATE FILE",
			params: map[string]string{
				"username":    user.Username,
				"usernameKey": "username",
				"fileKey":     "file",
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetFile(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(uploadResp, nil)

				querier.EXPECT().
					UpdateFile(
						gomock.Any(),
						gomock.Eq(db.UpdateFileParams{
							Username: user.Username,
							Data:     encoded,
						}),
					).
					Times(1).
					Return(uploadResp, nil)
			},
			setupAuth: func(
				t *testing.T, request *http.Request, tokenMaker *token.PasetoMaker,
			) {

				addAuthorization(
					t, request, tokenMaker, authorizationTypeToken, user.Username,
					time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchFile(t, recorder.Body, uploadResp)
			},
		},
		{
			name: "UNAUTHORIZED",
			params: map[string]string{
				"username":    "deidara",
				"usernameKey": "username",
				"fileKey":     "file",
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetFile(gomock.Any(), gomock.Any()).
					Times(0)

				querier.EXPECT().
					CreateFile(gomock.Any(), gomock.Any()).
					Times(0)

				querier.EXPECT().
					UpdateFile(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(
				t *testing.T, request *http.Request, tokenMaker *token.PasetoMaker,
			) {
				addAuthorization(
					t, request, tokenMaker, authorizationTypeToken, user.Username,
					time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "WRONG USERNAME KEY",
			params: map[string]string{
				"username":    "deidara",
				"usernameKey": "user",
				"fileKey":     "file",
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetFile(gomock.Any(), gomock.Any()).
					Times(0)

				querier.EXPECT().
					CreateFile(gomock.Any(), gomock.Any()).
					Times(0)

				querier.EXPECT().
					UpdateFile(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(
				t *testing.T, request *http.Request, tokenMaker *token.PasetoMaker,
			) {
				addAuthorization(
					t, request, tokenMaker, authorizationTypeToken, user.Username,
					time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "WRONG FILE KEY",
			params: map[string]string{
				"username":    "deidara",
				"usernameKey": "username",
				"fileKey":     "data",
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetFile(gomock.Any(), gomock.Any()).
					Times(0)

				querier.EXPECT().
					CreateFile(gomock.Any(), gomock.Any()).
					Times(0)

				querier.EXPECT().
					UpdateFile(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(
				t *testing.T, request *http.Request, tokenMaker *token.PasetoMaker,
			) {
				addAuthorization(
					t, request, tokenMaker, authorizationTypeToken, user.Username,
					time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "INTERNAL ERROR",
			params: map[string]string{
				"username":    user.Username,
				"usernameKey": "username",
				"fileKey":     "file",
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetFile(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(uploadResp, nil)

				querier.EXPECT().
					UpdateFile(
						gomock.Any(),
						gomock.Eq(db.UpdateFileParams{
							Username: user.Username,
							Data:     encoded,
						}),
					).
					Times(1).
					Return(db.File{}, sql.ErrConnDone)
			},
			setupAuth: func(
				t *testing.T, request *http.Request, tokenMaker *token.PasetoMaker,
			) {

				addAuthorization(
					t, request, tokenMaker, authorizationTypeToken, user.Username,
					time.Minute,
				)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			querier := mockdb.NewMockQuerier(ctrl)
			// build stubs
			tc.buildStubs(querier)

			// start test server and send request
			server := newTestServer(t, querier)
			recorder := httptest.NewRecorder()

			url := fmt.Sprint("/files/upload")
			buffer := bytes.Buffer{}
			mimeWriter := multipart.NewWriter(&buffer)
			err := mimeWriter.WriteField(
				tc.params["usernameKey"], tc.params["username"])
			require.NoError(t, err)
			formWriter, err := mimeWriter.CreateFormFile(
				tc.params["fileKey"], "test.csv")
			require.NoError(t, err)
			_, err = formWriter.Write([]byte(sampleCSV))
			require.NoError(t, err)
			mimeWriter.Close()

			request, err := http.NewRequest(http.MethodPost, url, &buffer)
			require.NoError(t, err)
			request.Header.Set("Content-Type", mimeWriter.FormDataContentType())

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchFile(t *testing.T, responseBody *bytes.Buffer, file db.File) {
	var serverResp fileResponse

	err := json.NewDecoder(responseBody).Decode(&serverResp)
	require.NoError(t, err)

	fmt.Println(serverResp.Error)
	require.Equal(t, file.ID, serverResp.File.ID)
	require.Equal(t, file.ChangedAt, serverResp.File.ChangedAt)
}
