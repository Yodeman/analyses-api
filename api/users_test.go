package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/yodeman/analyses-api/util"
	// "github.com/yodeman/analyses-api/token"
	mockdb "github.com/yodeman/analyses-api/dbase/mock"
	db "github.com/yodeman/analyses-api/dbase/sqlc"
)

type createUserParamsMatcher struct {
	createParams db.CreateUserParams
	password     string
}

func (params createUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(params.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	params.createParams.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(params.createParams, arg)
}

func (params createUserParamsMatcher) String() string {
	return fmt.Sprintf("matches database create user parameters.")
}

func TestCreateUser(t *testing.T) {
	user, password := randomUser(t)
	createUserParams := db.CreateUserParams{
		Username:       user.Username,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
	}

	req := createUserRequest{
		Username: user.Username,
		Password: password,
		Email:    user.Email,
	}

	testCases := []struct {
		name          string
		params        createUserRequest
		buildStubs    func(querier *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			params: req,
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					CreateUser(
						gomock.Any(),
						createUserParamsMatcher{
							createParams: createUserParams,
							password:     password,
						}).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "BAD REQUEST",
			params: createUserRequest{
				Username: user.Username,
				Email:    user.Email,
				Password: "123",
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "DUPLICATE",
			params: req,
			buildStubs: func(querier *mockdb.MockQuerier) {
				err := &pq.Error{
					Code: "23505", // postgress unique violation error code
				}
				querier.EXPECT().
					CreateUser(
						gomock.Any(),
						createUserParamsMatcher{
							createParams: createUserParams,
							password:     password,
						}).
					Times(1).
					Return(db.User{}, err)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name:   "INTERNAL ERROR",
			params: req,
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					CreateUser(
						gomock.Any(),
						createUserParamsMatcher{
							createParams: createUserParams,
							password:     password,
						}).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
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

			url := fmt.Sprint("/users/register")
			encodedParams, err := json.Marshal(tc.params)
			require.NoError(t, err)
			request, err := http.NewRequest(
				http.MethodPost, url, bytes.NewBuffer(encodedParams))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestLoginUser(t *testing.T) {
	user, password := randomUser(t)

	req := loginUserRequest{
		Username: user.Username,
		Password: password,
	}

	testCases := []struct {
		name          string
		params        loginUserRequest
		buildStubs    func(querier *mockdb.MockQuerier)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			params: req,
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetUser(gomock.Any(), req.Username).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "BAD REQUEST",
			params: loginUserRequest{
				Username: user.Username,
				Password: "123",
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "UNAUTHORIZED",
			params: loginUserRequest{
				Username: user.Username,
				Password: util.RandomPassword(),
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:   "INTERNAL ERROR",
			params: req,
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "NOT FOUND",
			params: req,
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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

			url := fmt.Sprint("/users/login")
			encodedParams, err := json.Marshal(tc.params)
			require.NoError(t, err)
			request, err := http.NewRequest(
				http.MethodGet, url, bytes.NewBuffer(encodedParams))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomPassword()
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomUser(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	return
}

func requireBodyMatchUser(t *testing.T, responseBody *bytes.Buffer, user db.User) {
	var serverResp userResponse

	err := json.NewDecoder(responseBody).Decode(&serverResp)
	require.NoError(t, err)

	require.Equal(t, user.Username, serverResp.User.Username)
	require.Equal(t, user.Email, serverResp.User.Email)
	require.Equal(t, user.PasswordChangedAt, serverResp.User.PasswordChangedAt)
	require.WithinDuration(t, user.CreatedAt, serverResp.User.CreatedAt, time.Second)
}
