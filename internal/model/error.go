package model

import "github.com/pkg/errors"

// ErrorNoteNotFound user not found error
var ErrorNoteNotFound = errors.New("user not found")

// ErrorFailToGenerateToken failed to generator token error
var ErrorFailToGenerateToken = errors.New("failed to generate token")

// ErrorInvalidRefereshToken invalid refresh token error
var ErrorInvalidRefereshToken = errors.New("invalid refresh token")

// ErrorMetadataNotProvided metadata is not provided
var ErrorMetadataNotProvided = errors.New("metadata is not provided")

// ErrorAuthorizationHeaderNotProvided authorization header is not provided error
var ErrorAuthorizationHeaderNotProvided = errors.New("authorization header is not provided")

// ErrorAuthorizationHeaderFormat invalid authorization header format
var ErrorAuthorizationHeaderFormat = errors.New("invalid authorization header format")

// ErrorAccessDenied access denied error
var ErrorAccessDenied = errors.New("access denied")
