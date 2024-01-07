package api

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

func TestLinearRegression(t *testing.T) {
	user, _ := randomUser(t)
	sampleCSV := util.RandomCSV(30, 10) // 30 rows and 10 cols csv file
	rows, cols, data, err := util.ParseCSVToFloatSlice(strings.NewReader(sampleCSV))
	require.NoError(t, err)
	matrix := mat.NewDense(rows, cols, data)
	byteData, err := matrix.MarshalBinary()
	require.NoError(t, err)
	encoded := base64.StdEncoding.EncodeToString(byteData)

	regReq := regressionRequest{
		Username: user.Username,
	}
	regResp := db.File{
		ID:       util.RandomInt(1, 1000),
		Username: user.Username,
		Data:     encoded,
	}

	testCases := []struct {
		name       string
		params     regressionRequest
		buildStubs func(querier *mockdb.MockQuerier)
		setupAuth  func(
			t *testing.T, request *http.Request, tokenMaker *token.PasetoMaker,
		)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			params: regReq,
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetFile(gomock.Any(), gomock.Eq(regReq.Username)).
					Times(1).
					Return(regResp, nil)
			},
			setupAuth: func(
				t *testing.T, request *http.Request, tokenMaker *token.PasetoMaker,
			) {
				addAuthorization(
					t, request, tokenMaker, authorizationTypeToken, user.Username,
					time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "BAD REQUEST",
			params: regressionRequest{Username: "1@2"},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetFile(gomock.Any(), gomock.Eq(regReq.Username)).
					Times(0)
			},
			setupAuth: func(
				t *testing.T, request *http.Request, tokenMaker *token.PasetoMaker,
			) {
				addAuthorization(
					t, request, tokenMaker, authorizationTypeToken, user.Username,
					time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "UNAUTHORIZED",
			params: regReq,
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetFile(gomock.Any(), gomock.Eq(regReq.Username)).
					Times(0)
			},
			setupAuth: func(
				t *testing.T, request *http.Request, tokenMaker *token.PasetoMaker,
			) {
				addAuthorization(
					t, request, tokenMaker, authorizationTypeToken, "12345",
					time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "INTERNAL ERROR",
			params: regReq,
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetFile(gomock.Any(), gomock.Eq(regReq.Username)).
					Times(1).
					Return(db.File{}, sql.ErrNoRows)
			},
			setupAuth: func(
				t *testing.T, request *http.Request, tokenMaker *token.PasetoMaker,
			) {
				addAuthorization(
					t, request, tokenMaker, authorizationTypeToken, user.Username,
					time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "CORRUPTED FILE",
			params: regReq,
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetFile(gomock.Any(), gomock.Eq(regReq.Username)).
					Times(1).
					Return(
						db.File{
							ID:       util.RandomInt(1, 1000),
							Username: user.Username,
							Data:     "ADIEDRYE=@$",
						},
						nil)
			},
			setupAuth: func(
				t *testing.T, request *http.Request, tokenMaker *token.PasetoMaker,
			) {
				addAuthorization(
					t, request, tokenMaker, authorizationTypeToken, user.Username,
					time.Minute)
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

			url := fmt.Sprint("/analyses/regression")
			encodedParams, err := json.Marshal(tc.params)
			require.NoError(t, err)
			request, err := http.NewRequest(
				http.MethodGet, url, bytes.NewBuffer(encodedParams))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
