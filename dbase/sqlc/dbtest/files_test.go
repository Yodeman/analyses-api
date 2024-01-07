package dbtest

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	mockdb "github.com/yodeman/analyses-api/dbase/mock"
	db "github.com/yodeman/analyses-api/dbase/sqlc"
	"github.com/yodeman/analyses-api/util"
)

func TestCreateFile(t *testing.T) {
	user, _ := randomUser(t)
	data, err := util.RandomData()
	require.NoError(t, err)

	createFileParams := db.CreateFileParams{
		Username: user.Username,
		Data:     data,
	}

	file := db.File{
		Username: user.Username,
		Data:     data,
	}

	var ctx context.Context

	testCases := []struct {
		name        string
		buildStubs  func(querier *mockdb.MockQuerier)
		checkResult func(t *testing.T, result db.File, err error)
	}{
		{
			name: "OK",
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					CreateFile(gomock.Any(), gomock.Eq(createFileParams)).
					Times(1).
					Return(file, nil)
			},
			checkResult: func(t *testing.T, result db.File, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, result)
				require.Equal(t, file.Username, result.Username)
				require.Equal(t, file.Data, result.Data)
				require.WithinDuration(t, user.CreatedAt, result.CreatedAt, time.Second)
			},
		},
		{
			name: "INTERNAL ERROR",
			buildStubs: func(querier *mockdb.MockQuerier) {

				querier.EXPECT().
					CreateFile(gomock.Any(), gomock.Eq(createFileParams)).
					Times(1).
					Return(db.File{}, sql.ErrConnDone)
			},
			checkResult: func(t *testing.T, result db.File, err error) {
				require.Error(t, err)
				require.Empty(t, result)
				require.Zero(t, result.Username)
				require.Zero(t, result.Data)
				require.Zero(t, result.CreatedAt)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			testQuerier := mockdb.NewMockQuerier(ctrl)

			//build stubs
			tc.buildStubs(testQuerier)

			result, err := testQuerier.CreateFile(ctx, createFileParams)

			tc.checkResult(t, result, err)
		})
	}
}

func TestGetFile(t *testing.T) {
	user, _ := randomUser(t)
	data, err := util.RandomData()
	require.NoError(t, err)

	file := db.File{
		Username: user.Username,
		Data:     data,
	}

	var ctx context.Context

	testCases := []struct {
		name        string
		param       string
		buildStubs  func(querier *mockdb.MockQuerier)
		checkResult func(t *testing.T, result db.File, err error)
	}{
		{
			name:  "OK",
			param: user.Username,
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					GetFile(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(file, nil)
			},
			checkResult: func(t *testing.T, result db.File, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, result)
				require.Equal(t, file.Username, result.Username)
				require.Equal(t, file.Data, result.Data)
				require.WithinDuration(t, user.CreatedAt, result.CreatedAt, time.Second)
			},
		},
		{
			name:  "INTERNAL ERROR",
			param: user.Username,
			buildStubs: func(querier *mockdb.MockQuerier) {

				querier.EXPECT().
					GetFile(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(db.File{}, sql.ErrConnDone)
			},
			checkResult: func(t *testing.T, result db.File, err error) {
				require.Error(t, err)
				require.Empty(t, result)
				require.Zero(t, result.Username)
				require.Zero(t, result.Data)
				require.Zero(t, result.CreatedAt)
			},
		},
		{
			name:  "NOT FOUND",
			param: "",
			buildStubs: func(querier *mockdb.MockQuerier) {

				querier.EXPECT().
					GetFile(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.File{}, sql.ErrNoRows)
			},
			checkResult: func(t *testing.T, result db.File, err error) {
				require.Error(t, err)
				require.Empty(t, result)
				require.Zero(t, result.Username)
				require.Zero(t, result.Data)
				require.Zero(t, result.CreatedAt)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			testQuerier := mockdb.NewMockQuerier(ctrl)

			//build stubs
			tc.buildStubs(testQuerier)

			result, err := testQuerier.GetFile(ctx, tc.param)

			tc.checkResult(t, result, err)
		})
	}
}

func TestUpdateFile(t *testing.T) {
	user, _ := randomUser(t)
	data, err := util.RandomData()
	require.NoError(t, err)

	updateFileParams := db.UpdateFileParams{
		Username: user.Username,
		Data:     data,
	}

	file := db.File{
		Username: user.Username,
		Data:     data,
	}

	var ctx context.Context

	testCases := []struct {
		name        string
		param       db.UpdateFileParams
		buildStubs  func(querier *mockdb.MockQuerier)
		checkResult func(t *testing.T, result db.File, err error)
	}{
		{
			name:  "OK",
			param: updateFileParams,
			buildStubs: func(querier *mockdb.MockQuerier) {
				querier.EXPECT().
					UpdateFile(gomock.Any(), gomock.Eq(updateFileParams)).
					Times(1).
					Return(file, nil)
			},
			checkResult: func(t *testing.T, result db.File, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, result)
				require.Equal(t, file.Username, result.Username)
				require.Equal(t, file.Data, result.Data)
				require.WithinDuration(t, user.CreatedAt, result.CreatedAt, time.Second)
			},
		},
		{
			name:  "INTERNAL ERROR",
			param: updateFileParams,
			buildStubs: func(querier *mockdb.MockQuerier) {

				querier.EXPECT().
					UpdateFile(gomock.Any(), gomock.Eq(updateFileParams)).
					Times(1).
					Return(db.File{}, sql.ErrConnDone)
			},
			checkResult: func(t *testing.T, result db.File, err error) {
				require.Error(t, err)
				require.Empty(t, result)
				require.Zero(t, result.Username)
				require.Zero(t, result.Data)
				require.Zero(t, result.CreatedAt)
			},
		},
		{
			name: "NOT FOUND",
			param: db.UpdateFileParams{
				Username: "",
				Data:     data,
			},
			buildStubs: func(querier *mockdb.MockQuerier) {
				param := db.UpdateFileParams{
					Username: "",
					Data:     data,
				}
				querier.EXPECT().
					UpdateFile(gomock.Any(), gomock.Eq(param)).
					Times(1).
					Return(db.File{}, sql.ErrNoRows)
			},
			checkResult: func(t *testing.T, result db.File, err error) {
				require.Error(t, err)
				require.Empty(t, result)
				require.Zero(t, result.Username)
				require.Zero(t, result.Data)
				require.Zero(t, result.CreatedAt)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			testQuerier := mockdb.NewMockQuerier(ctrl)

			//build stubs
			tc.buildStubs(testQuerier)

			result, err := testQuerier.UpdateFile(ctx, tc.param)

			tc.checkResult(t, result, err)
		})
	}
}
