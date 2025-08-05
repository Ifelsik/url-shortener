package app

import (
	"github.com/Ifelsik/url-shortener/internal/app/url"
	"github.com/Ifelsik/url-shortener/internal/app/user"
)

type URLService struct {
	AddURL     url.AddURL
	GetByShort url.GetURLByShort
}

type UserService struct {
	AddUser user.AddUser
}
