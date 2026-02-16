package store

import "errors"

var ErrTaskAlreadyActive = errors.New("task already active")
var ErrTaskNotActive = errors.New("task not active")
