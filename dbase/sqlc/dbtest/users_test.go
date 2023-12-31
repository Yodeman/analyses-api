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


func TestCreateUser(t *testing.T) {
    user, password := randomUser(t)
    createUserParams := db.CreateUserParams{
        Username:       user.Username,
        HashedPassword: user.HashedPassword,
        Email:          user.Email,
    }
    var ctx context.Context

    testCases := []struct{
        name        string
        buildStubs  func(querier *mockdb.MockQuerier)
        checkResult func(t *testing.T, result db.User, err error)
    }{
        {
            name:       "OK",
            buildStubs: func(querier *mockdb.MockQuerier) {
                querier.EXPECT().
                    CreateUser(gomock.Any(), gomock.Eq(createUserParams)).
                    Times(1).
                    Return(user, nil)
            },
            checkResult:    func(t *testing.T, result db.User, err error) {
                require.NoError(t, err)
                require.NotEmpty(t, result)
                require.Equal(t, user.Username, result.Username)
                require.Equal(t, user.Email, result.Email)
                require.WithinDuration(t, user.CreatedAt, result.CreatedAt, time.Second)
                require.NoError(t, util.CheckPassword(password, result.HashedPassword))
            },
        },
        {
            name:       "INTERNAL ERROR",
            buildStubs: func(querier *mockdb.MockQuerier) {

                querier.EXPECT().
                    CreateUser(gomock.Any(), gomock.Eq(createUserParams)).
                    Times(1).
                    Return(db.User{}, sql.ErrConnDone)
            },
            checkResult:    func(t *testing.T, result db.User, err error) {
                require.Error(t, err)
                require.Empty(t, result)
                require.Zero(t, result.Username)
                require.Zero(t, result.Email)
                require.Zero(t, result.CreatedAt)
                require.Error(t, util.CheckPassword(password, result.HashedPassword))
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

            result, err := testQuerier.CreateUser(ctx, createUserParams)

            tc.checkResult(t, result, err)
        })
    }
}

func TestGetUser(t *testing.T) {
    user, password := randomUser(t)
    var ctx context.Context

    testCases := []struct{
        name        string
        param       string
        buildStubs  func(querier *mockdb.MockQuerier)
        checkResult func(t *testing.T, result db.User, err error)
    }{
        {
            name:       "OK",
            param:      user.Username,
            buildStubs: func(querier *mockdb.MockQuerier) {
                querier.EXPECT().
                    GetUser(gomock.Any(), gomock.Eq(user.Username)).
                    Times(1).
                    Return(user, nil)
            },
            checkResult:    func(t *testing.T, result db.User, err error) {
                require.NoError(t, err)
                require.NotEmpty(t, result)
                require.Equal(t, user.Username, result.Username)
                require.Equal(t, user.Email, result.Email)
                require.WithinDuration(t, user.CreatedAt, result.CreatedAt, time.Second)
                require.NoError(t, util.CheckPassword(password, result.HashedPassword))
            },
        },
        {
            name:       "NOT FOUND",
            param:      "",
            buildStubs: func(querier *mockdb.MockQuerier) {

                querier.EXPECT().
                    GetUser(gomock.Any(), gomock.Any()).
                    Times(1).
                    Return(db.User{}, sql.ErrNoRows)
            },
            checkResult:    func(t *testing.T, result db.User, err error) {
                require.Error(t, err)
                require.Empty(t, result)
                require.Zero(t, result.Username)
                require.Zero(t, result.Email)
                require.Zero(t, result.CreatedAt)
                require.Error(t, util.CheckPassword(password, result.HashedPassword))
            },
        },
        {
            name:       "INTERNAL ERROR",
            param:      user.Username,
            buildStubs: func(querier *mockdb.MockQuerier) {

                querier.EXPECT().
                    GetUser(gomock.Any(), gomock.Eq(user.Username)).
                    Times(1).
                    Return(db.User{}, sql.ErrConnDone)
            },
            checkResult:    func(t *testing.T, result db.User, err error) {
                require.Error(t, err)
                require.Empty(t, result)
                require.Zero(t, result.Username)
                require.Zero(t, result.Email)
                require.Zero(t, result.CreatedAt)
                require.Error(t, util.CheckPassword(password, result.HashedPassword))
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

            result, err := testQuerier.GetUser(ctx, tc.param)

            tc.checkResult(t, result, err)
        })
    }
}


func randomUser(t *testing.T) (user db.User, password string) {
    password = util.RandomPassword()
    hashedPassword, err := util.HashPassword(password)
    require.NoError(t, err)

    user = db.User{
        Username:       util.RandomUser(),
        HashedPassword: hashedPassword,
        Email:          util.RandomEmail(),
    }

    return
}
