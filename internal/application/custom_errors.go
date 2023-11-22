package application

import "errors"

var ErrLoginAlreadyInUse = errors.New("login already in use")
var ErrLoginOrPasswordIsInvalid = errors.New("login or password is invalid")
var ErrNotValidNumber = errors.New("order number is invalid")
var ErrUploadedByThisUser = errors.New("order was uploaded by this user")
var ErrUploadedByAnotherUser = errors.New("order was uploaded by another user")
