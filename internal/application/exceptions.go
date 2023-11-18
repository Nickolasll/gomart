package application

import "errors"

var ErrLoginAlreadyInUse = errors.New("login already in use")
var ErrLoginOrPasswordIsInvalid = errors.New("login or password is invalid")
var ErrNotValidNumber = errors.New("order number is invalid")
var ErrOrderUploadedByThisUser = errors.New("order was uploaded by this user")
var ErrOrderUploadedByAnotherUser = errors.New("order was uploaded by another user")
