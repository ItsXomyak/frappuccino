package cerrors

import (
	"errors"
	"log/slog"
)

var (
	ErrNotExist         = errors.New("not exist")
	ErrExist            = errors.New("exist")
	ErrNameIsNotValid   = errors.New("the name is not valid")
	ErrIsNotEmpty       = errors.New("is not empty")
	ErrOrderNotFound    = errors.New("order id not found")
	ErrMenuItemNotFound = errors.New("menu item id not found")
)

func NotExist() error {
	slog.Error("Error occurred", "details", ErrNotExist.Error())
	return ErrNotExist
}

func Exist() error {
	slog.Error("Error occurred", "details", ErrExist.Error())
	return ErrExist
}

func NameIsNotValid() error {
	slog.Error("Error occurred", "details", ErrNameIsNotValid.Error())
	return ErrNameIsNotValid
}

func IsNotEmpty() error {
	slog.Error("Error occurred", "details", ErrIsNotEmpty.Error())
	return ErrIsNotEmpty
}

func OrderNotFound() error {
	slog.Error("Error occurred", "details", ErrOrderNotFound.Error())
	return ErrOrderNotFound
}
