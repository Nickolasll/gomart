package domain

import "errors"

var ErrInsufficientFunds = errors.New("insufficient funds on current user balance")
var ErrDocumentNotFound = errors.New("requested document not found")
var ErrAccrualIsBusy = errors.New("accrual service is not ready")
