package url_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Ifelsik/url-shortener/internal/app/url"
	urldom "github.com/Ifelsik/url-shortener/internal/domain/url"
	"github.com/Ifelsik/url-shortener/internal/mocks"
	"github.com/stretchr/testify/suite"
)

type GetURLByShortTestSuite struct {
	suite.Suite

	// dependencies
	mockURLRepo   *mocks.MockURLRepository
	mockValidator *mocks.MockValidator
}

func (s *GetURLByShortTestSuite) SetupTest() {
	s.mockURLRepo = mocks.NewMockURLRepository(s.T())
	s.mockValidator = mocks.NewMockValidator(s.T())
}

func (s *GetURLByShortTestSuite) TestSuccess() {
	s.T().Log("Should return original url with short key if success")

	// Arrange
	const shortKey = "shortedYa"
	const originalURL = "https://ya.ru"

	request := &url.GetURLByShortRequest{
		ShortKey: shortKey,
	}

	s.mockValidator.EXPECT().ValidateStruct(request).Return(nil)
	s.mockURLRepo.EXPECT().GetByShortKey(context.TODO(), shortKey).
		Return(&urldom.URL{ShortKey: shortKey, OriginalURL: originalURL}, nil)

	want := &url.GetURLByShortResponse{
		ShortURL:    shortKey,
		OriginalURL: originalURL,
	}

	testedHandler := url.NewGetURLByShortKey(s.mockURLRepo, s.mockValidator)

	// Act
	got, err := testedHandler.Handle(context.TODO(), request)

	// Assert
	s.Require().NoError(err)
	s.Require().Equal(want, got)
}

func (s *GetURLByShortTestSuite) TestNoShortKeyExists() {
	s.T().Log("Should return error if no short key exists")

	// Arrange
	const shortKey = "shortedYa"

	request := &url.GetURLByShortRequest{
		ShortKey: shortKey,
	}

	s.mockValidator.EXPECT().ValidateStruct(request).Return(nil)
	s.mockURLRepo.EXPECT().GetByShortKey(context.TODO(), shortKey).
		Return(nil, urldom.ErrNoURL)

	testedHandler := url.NewGetURLByShortKey(s.mockURLRepo, s.mockValidator)

	// Act
	got, err := testedHandler.Handle(context.TODO(), request)

	// Assert
	s.Require().Error(err)
	s.Require().ErrorIs(err, url.ErrNotFound)
	s.Require().Nil(got)
}

func (s *GetURLByShortTestSuite) TestUnknownDomainError() {
	s.T().Log("Should return error if got unknown domain layer error")

	// Arrange
	const shortKey = "shortedYa"

	request := &url.GetURLByShortRequest{
		ShortKey: shortKey,
	}

	s.mockValidator.EXPECT().ValidateStruct(request).Return(nil)
	s.mockURLRepo.EXPECT().GetByShortKey(context.TODO(), shortKey).
		Return(nil, errors.New("unknown error"))

	testedHandler := url.NewGetURLByShortKey(s.mockURLRepo, s.mockValidator)

	// Act
	got, err := testedHandler.Handle(context.TODO(), request)

	// Assert
	s.Require().Error(err)
	s.Require().NotErrorIs(err, url.ErrNotFound)
	s.Require().Nil(got)
}

func (s *GetURLByShortTestSuite) TestInvalidRequest() {
	s.T().Log("Should return error if request validation fails")

	// Arrange
	const shortKey = "shortedYa"

	request := &url.GetURLByShortRequest{
		ShortKey: shortKey,
	}

	s.mockValidator.EXPECT().ValidateStruct(request).Return(errors.New("validation error"))

	testedHandler := url.NewGetURLByShortKey(s.mockURLRepo, s.mockValidator)

	// Act
	got, err := testedHandler.Handle(context.TODO(), request)

	// Assert
	s.Require().Error(err)
	s.Require().Nil(got)
}

func (s *GetURLByShortTestSuite) TestNilRequest() {
	s.T().Log("Should return error if request is nil")

	// Arrange
	testedHandler := url.NewGetURLByShortKey(s.mockURLRepo, s.mockValidator)

	ctx := context.TODO()
	var request *url.GetURLByShortRequest = nil

	// Act
	got, err := testedHandler.Handle(ctx, request)

	// Assert
	s.Require().ErrorIs(err, url.ErrEmptyRequest)
	s.Require().Nil(got)
}

func TestGetURLByShortTestSuite(t *testing.T) {
	suite.Run(t, new(GetURLByShortTestSuite))
}
