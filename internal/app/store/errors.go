package store

import "errors"

var (
	ErrNotAuth = errors.New("Not auth")
	ErrWrongReq = errors.New("Wrong req")
	ErrRecordNotFound = errors.New("Record not found")
	ErrNotValidToken = errors.New("Not valid token")
	ErrNotCurrectRoomName = errors.New("Not valid room name")
	ErrWrongRequset = errors.New("Wrong request")
	InternalError = errors.New("Internal error")
	ErrBigFile = errors.New("Upload file limit")
	ErrIvalidFieFile = errors.New("InvalidFile")
)
