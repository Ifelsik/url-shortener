package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Ifelsik/url-shortener/internal/app/user"
	domain "github.com/Ifelsik/url-shortener/internal/domain/user"
	"github.com/Ifelsik/url-shortener/internal/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_addUser_Handle(t *testing.T) {
	const expectedToken = "a1b2c3d4-e5f6-7890-1234-567890abcdef"

	type args struct {
		ctx     context.Context
		request *user.AddUserRequest
	}
	tests := []struct {
		name    string
		arrange func(t *testing.T) *user.AddUserProvider
		args    args
		want    *user.AddUserResponse
		wantErr bool
	}{
		{
			name: "Should add new user and return user token",
			//nolint:thelper
			arrange: func(t *testing.T) *user.AddUserProvider {
				nowTimeMocked := time.Date(2025, 1, 25, 0, 0, 0, 0, time.UTC)
				expectedExpiresAtMocked := nowTimeMocked.Add(30 * 24 * time.Hour)

				timeMock := mocks.NewMockTiming(t)
				timeMock.EXPECT().Now().Return(nowTimeMocked)
				timeMock.EXPECT().
					AfterNow(30 * 24 * time.Hour).
					Return(expectedExpiresAtMocked)

				identifierMock := mocks.NewMockIdentifier(t)
				identifierMock.EXPECT().String().Return(expectedToken)

				userRepositoryMock := mocks.NewMockUserRepository(t)
				userRepositoryMock.EXPECT().GetByToken(mock.Anything, "").
					Return(nil, domain.ErrNoUser)

				userRepositoryMock.EXPECT().
					Add(mock.Anything,
						mock.MatchedBy(func(user *domain.User) bool {
							require.Equal(t, expectedToken, user.Token)
							require.Equal(t, nowTimeMocked, user.CreatedAt)
							require.Equal(t, expectedExpiresAtMocked, user.ExpiresAt)

							return true
						}),
					).Return(&domain.User{Token: expectedToken}, nil)

				return user.NewAddUser(
					userRepositoryMock,
					timeMock,
					identifierMock,
				)
			},
			args: args{
				ctx: context.TODO(),
				request: &user.AddUserRequest{
					UserToken: "",
				},
			},
			want: &user.AddUserResponse{
				UserToken: expectedToken,
			},
			wantErr: false,
		},
		{
			name: "Should return already user's token if provided user exists",
			//nolint:thelper
			arrange: func(t *testing.T) *user.AddUserProvider {
				userRepositoryMock := mocks.NewMockUserRepository(t)

				userRepositoryMock.EXPECT().GetByToken(mock.Anything, expectedToken).
					Return(&domain.User{Token: expectedToken}, nil)

				return user.NewAddUser(
					userRepositoryMock,
					nil,
					nil,
				)
			},
			args: args{
				ctx: context.TODO(),
				request: &user.AddUserRequest{
					UserToken: expectedToken,
				},
			},
			want: &user.AddUserResponse{
				UserToken: expectedToken,
			},
			wantErr: false,
		},
		{
			name: "Should return error if adding user fails",
			//nolint:thelper
			arrange: func(t *testing.T) *user.AddUserProvider {
				userRepositoryMock := mocks.NewMockUserRepository(t)
				timingMock := mocks.NewMockTiming(t)
				identifierMock := mocks.NewMockIdentifier(t)

				var dummyTime time.Time

				userRepositoryMock.EXPECT().GetByToken(mock.Anything, expectedToken).
					Return(nil, domain.ErrNoUser)
				userRepositoryMock.EXPECT().Add(mock.Anything, mock.Anything).
					Return(nil, errors.New("some error"))

				identifierMock.EXPECT().String().Return(expectedToken)

				timingMock.EXPECT().Now().Return(dummyTime)
				timingMock.EXPECT().AfterNow(30 * 24 * time.Hour).Return(dummyTime)

				return user.NewAddUser(
					userRepositoryMock,
					timingMock,
					identifierMock,
				)
			},
			args: args{
				ctx: context.TODO(),
				request: &user.AddUserRequest{
					UserToken: expectedToken,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addUserProvider := tt.arrange(t)

			require.NotNil(t, addUserProvider,
				"addUserProvider.Handle() is nil", tt.name)

			got, err := addUserProvider.Handle(tt.args.ctx, tt.args.request)

			require.Equal(t, tt.wantErr, err != nil)

			require.Equal(t, tt.want, got)
		})
	}
}
