package domain

import "errors"

var (
	// ErrInsufficientBalance は残高不足（不変条件違反）を表す
	ErrInsufficientBalance = errors.New("domain: insufficient balance")
	// ErrInvalidAmount は送金額が0以下であることを表す
	ErrInvalidAmount = errors.New("domain: amount must be positive")
	// ErrAccountNotFound は口座が存在しないことを表す
	ErrAccountNotFound = errors.New("domain: account not found")
)
