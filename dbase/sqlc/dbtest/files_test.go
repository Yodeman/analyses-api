package dbtest

import (
    "context"
    "database/sql"
    "testing"
    "time"

    "github.com/stretchr/testify/require"
    gomock "go.uber.org/mock/gomock"

    "github.com/yodeman/analyses-api/util"
    db "github.com/yodeman/analyses-api/dbase/sqlc"
    mockdb "github.com/yodeman/analyses-api/dbase/mock"
)

func TestCreateFile(t *testing.T) {
    user, _ := randomUser(t)
    data, err := util.RandomData()
    require.NoError(t, err)

    createFileParams := db.CreateFileParams{
        Username:   user.Username,
        Data:       string(data),
    }

    file := db.File{
        Username:   user.Username,
        Data:       string(data),
    }

    var ctx context.Context

    testCases := []struct{
        name        string
        buildStubs  func(querier *mockdb.MockQuerier)
        checkResult func(t *testing.T, result db.File, err error)
    }{
        {
            name:       "OK",
            buildStubs: func(querier *mockdb.MockQuerier) {
                querier.EXPECT().
                    CreateFile(gomock.Any(), gomock.Eq(createFileParams)).
                    Times(1).
                    Return(file, nil)
            },
            checkResult:    func(t *testing.T, result db.File, err error) {
                require.NoError(t, err)
                require.NotEmpty(t, result)
                require.Equal(t, file.Username, result.Username)
                require.Equal(t, file.Data, result.Data)
                require.WithinDuration(t, user.CreatedAt, result.CreatedAt, time.Second)
            },
        },
        {
            name:       "INTERNAL ERROR",
            buildStubs: func(querier *mockdb.MockQuerier) {

                querier.EXPECT().
                    CreateFile(gomock.Any(), gomock.Eq(createFileParams)).
                    Times(1).
                    Return(db.File{}, sql.ErrConnDone)
            },
            checkResult:    func(t *testing.T, result db.File, err error) {
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

        t.Run(tc.name, func(t *testing.T){
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
