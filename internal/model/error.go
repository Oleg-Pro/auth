package model

import "github.com/pkg/errors"

// ErrorNoteNotFound user not found error
var ErrorNoteNotFound = errors.New("user not found")

// ErrorFailToGenerateToken failed to generator token error
var ErrorFailToGenerateToken = errors.New("failed to generate token")
