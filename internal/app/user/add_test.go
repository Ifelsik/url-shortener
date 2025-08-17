package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/Ifelsik/url-shortener/internal/app/user"
	domain "github.com/Ifelsik/url-shortener/internal/domain/user"
	"github.com/Ifelsik/url-shortener/internal/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_addUser_Handle(t *testing.T) {
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
		// TODO: Add test cases.
		{
			name: "Should add new user and return user token",
			arrange: func(t *testing.T) *user.AddUserProvider {
				nowTimeMocked := time.Date(2025, 1, 25, 0, 0, 0, 0, time.UTC)
				expectToken := "a1b2c3d4-e5f6-7890-1234-567890abcdef"
				expectedExpiresAtMocked := nowTimeMocked.Add(30 * 24 * time.Hour)

				timeMock := mocks.NewMockTiming(t)
				timeMock.EXPECT().Now().Return(nowTimeMocked)
				timeMock.EXPECT().
					AfterNow(30 * 24 * time.Hour).
					Return(expectedExpiresAtMocked)

				identifierMock := mocks.NewMockIdentifier(t)
				identifierMock.EXPECT().String().Return(expectToken)

				userRepositoryMock := mocks.NewMockUserRepository(t)
				userRepositoryMock.EXPECT().GetByToken(mock.Anything, "").
					Return(nil, domain.ErrNoUser)

				userRepositoryMock.EXPECT().
					Add(mock.Anything,
						mock.MatchedBy(func(user *domain.User) bool {
							require.Equal(t, expectToken, user.Token)
							require.Equal(t, nowTimeMocked, user.CreatedAt)
							require.Equal(t, expectedExpiresAtMocked, user.ExpiresAt)

							return true
						}),
					).Return(&domain.User{Token: expectToken}, nil)

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
				UserToken: "a1b2c3d4-e5f6-7890-1234-567890abcdef",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addUserProvider := tt.arrange(t)

			require.NotNilf(t, addUserProvider,
				"addUserProvider.Handle() is nil", tt.name)

			got, err := addUserProvider.Handle(tt.args.ctx, tt.args.request)

			require.Equal(t, tt.wantErr, err != nil)

			require.Equal(t, tt.want, got)
		})
	}
}
