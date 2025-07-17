package service

import "fmt"

type ErrRegisteredEmail struct {
	Email string
}

func (err *ErrRegisteredEmail) Error() string {
	return fmt.Sprintf("Registered Email: email=%s", err.Email)
}

type ErrRegisteredHandleName struct {
	HandleName string
}

func (err *ErrRegisteredHandleName) Error() string {
	return fmt.Sprintf("Registered HandleName: handlename=%s", err.HandleName)
}

type ErrUserNotFound struct {
	UserId uint
	Email  string
}

func (err *ErrUserNotFound) Error() string {
	return fmt.Sprintf("User Not Found: user_id=%v, email=%v", err.UserId, err.Email)
}

type ErrCredentialNotFound struct {
	UserId uint
}

func (err *ErrCredentialNotFound) Error() string {
	return fmt.Sprintf("Credential Not Found: user_id=%v", err.UserId)
}

type ErrPasswordMissMatch struct {
	UserId uint
}

func (err *ErrPasswordMissMatch) Error() string {
	return fmt.Sprintf("Password MissMatch with HashedOne: user_id=%v", err.UserId)
}

type ErrInvalidRefreshTokenProvided struct {
}

func (err *ErrInvalidRefreshTokenProvided) Error() string {
	return fmt.Sprintf("Invalid RefreshToken Provided")
}

type ErrExpiredRefreshTokenProvided struct {
}

func (err *ErrExpiredRefreshTokenProvided) Error() string {
	return fmt.Sprintf("Expired RefreshToken Provided")
}
