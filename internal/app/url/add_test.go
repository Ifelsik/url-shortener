package url_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Ifelsik/url-shortener/internal/app/url"
	"github.com/Ifelsik/url-shortener/internal/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	urldom "github.com/Ifelsik/url-shortener/internal/domain/url"
	userdom "github.com/Ifelsik/url-shortener/internal/domain/user"
)

type AddURLTestSuite struct {
	suite.Suite

	// dependencies
	mockURLRepo        *mocks.MockURLRepository
	mockUserRepo       *mocks.MockUserRepository
	mockTimingProvider *mocks.MockTiming
	mockBase62Provider *mocks.MockBase62Provider
	mockHasher         *mocks.MockHasher
	mockValidator      *mocks.MockValidator
}

func (s *AddURLTestSuite) SetupTest() {
	t := s.T()

	s.mockURLRepo = mocks.NewMockURLRepository(t)
	s.mockUserRepo = mocks.NewMockUserRepository(t)
	s.mockTimingProvider = mocks.NewMockTiming(t)
	s.mockBase62Provider = mocks.NewMockBase62Provider(t)
	s.mockHasher = mocks.NewMockHasher(t)
	s.mockValidator = mocks.NewMockValidator(t)
}

func (s *AddURLTestSuite) TestNilRequest() {
	s.T().Log("Should return error if request is nil")

	// Arrange
	testedAddURL := url.NewAddURL(
		s.mockURLRepo,
		s.mockUserRepo,
		s.mockTimingProvider,
		s.mockBase62Provider,
		s.mockHasher,
		s.mockValidator,
	)

	ctx := context.TODO()
	var request *url.AddURLRequest = nil

	// Act
	got, err := testedAddURL.Handle(ctx, request)

	// Assert
	s.Require().ErrorIs(err, url.ErrEmptyRequest)
	s.Require().Nil(got)
}

func (s *AddURLTestSuite) TestInvalidRequest() {
	s.T().Log("Should return error if request is invalid")

	// Arrange
	ctx := context.TODO()
	request := &url.AddURLRequest{
		OriginalURL: "invalid url",
	}

	s.mockValidator.EXPECT().ValidateStruct(request).
		Return(errors.New("incalid request"))

	testedAddURL := url.NewAddURL(
		s.mockURLRepo,
		s.mockUserRepo,
		s.mockTimingProvider,
		s.mockBase62Provider,
		s.mockHasher,
		s.mockValidator,
	)

	// Act
	got, err := testedAddURL.Handle(ctx, request)

	// Assert
	s.Require().Error(err)
	s.Require().Nil(got)
}

func (s *AddURLTestSuite) TestInvalidUserToken() {
	s.T().Log("Should return error if user token is invalid")

	// Arrange
	const invalidToken = "invalid token"

	ctx := context.TODO()
	request := &url.AddURLRequest{
		UserToken:   invalidToken,
		OriginalURL: "https://google.com",
	}

	s.mockValidator.EXPECT().ValidateStruct(request).Return(nil)
	s.mockUserRepo.EXPECT().GetByToken(ctx, invalidToken).Return(nil, userdom.ErrNoUser)

	testedAddURL := url.NewAddURL(
		s.mockURLRepo,
		s.mockUserRepo,
		s.mockTimingProvider,
		s.mockBase62Provider,
		s.mockHasher,
		s.mockValidator,
	)

	// Act
	got, err := testedAddURL.Handle(ctx, request)

	// Assert
	s.Require().Error(err)
	s.Require().Nil(got)
}

func (s *AddURLTestSuite) TestUnableToSaveURL() {
	s.T().Log("Should return error if unable to save url")

	// Arrange
	const URL = "https://ya.ru"
	const userToken = "token"

	ctx := context.TODO()
	request := &url.AddURLRequest{
		OriginalURL: URL,
		UserToken:   userToken,
	}

	s.mockValidator.EXPECT().ValidateStruct(request).Return(nil)

	const ID = 1
	s.mockUserRepo.EXPECT().GetByToken(ctx, userToken).
		Return(&userdom.User{ID: ID}, nil)

	const hash = "hashedURL"
	s.mockHasher.EXPECT().String(URL).Return(hash)

	const shortKey = "shortKey"
	s.mockBase62Provider.EXPECT().EncodeToString([]byte(hash)).Return(shortKey)

	s.mockTimingProvider.EXPECT().Now().Return(time.Time{})

	s.mockURLRepo.EXPECT().Add(ctx, mock.Anything).
		Return(nil, errors.New("failed to add"))
	
	testedAddURL := url.NewAddURL(
		s.mockURLRepo,
		s.mockUserRepo,
		s.mockTimingProvider,
		s.mockBase62Provider,
		s.mockHasher,
		s.mockValidator,
	)

	// Act
	got, err := testedAddURL.Handle(ctx, request)

	// Assert
	s.Require().Error(err)
	s.Require().Nil(got)
}

func (s *AddURLTestSuite) TestSuccess() {
	s.T().Log("Should return struct with original and short url if success")

	// Arrange
	const URL = "https://google.ru"
	const userToken = "token"

	ctx := context.TODO()
	request := &url.AddURLRequest{
		OriginalURL: URL,
		UserToken:   userToken,
	}

	s.mockValidator.EXPECT().ValidateStruct(request).Return(nil)

	const ID = uint64(1)
	s.mockUserRepo.EXPECT().GetByToken(ctx, userToken).
		Return(&userdom.User{ID: ID}, nil)

	const hash = "hashedURL"
	s.mockHasher.EXPECT().String(URL).Return(hash)

	const shortKey = "shortKey"
	s.mockBase62Provider.EXPECT().EncodeToString([]byte(hash)).Return(shortKey)

	s.mockTimingProvider.EXPECT().Now().Return(time.Time{})

	s.mockURLRepo.EXPECT().Add(ctx, mock.MatchedBy(func (url *urldom.URL) bool {
		s.Equal(URL, url.OriginalURL)
		s.Equal(ID, url.User)
		s.Equal(shortKey, url.ShortKey)
		return true
	})).
		Return(&urldom.URL{OriginalURL: URL, ShortKey: shortKey}, nil)
	
	testedAddURL := url.NewAddURL(
		s.mockURLRepo,
		s.mockUserRepo,
		s.mockTimingProvider,
		s.mockBase62Provider,
		s.mockHasher,
		s.mockValidator,
	)

	expected := &url.AddURLResponse{
		OriginalURL: URL,
		ShortURL:    shortKey,
	}

	// Act
	got, err := testedAddURL.Handle(ctx, request)

	// Assert
	s.Require().NoError(err)
	s.Require().Equal(expected, got)
}

func TestAddURLTestSuite(t *testing.T) {
	suite.Run(t, new(AddURLTestSuite))
}
