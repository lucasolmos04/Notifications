package service

import "errors"

var ErrNotFound = errors.New("resource not found")
var ErrUnsupportedChannel = errors.New("unsupported channel")
