package repository

import "errors"

var ErrNoID = errors.New("an ID needs to be supplied")
var ErrBodyTooBig = errors.New("body supplied exceeded the maximum size of 5 MB")
var ErrNoAuthHeader = errors.New("authorization headers must be supplied")
