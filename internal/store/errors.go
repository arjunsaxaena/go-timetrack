package store

import "errors"

var ErrTaskAlreadyActive = errors.New("task already active")
var ErrTaskNotActive = errors.New("task not active")
var ErrLogNotFound = errors.New("log not found")
var ErrInvalidTimeRange = errors.New("start time cannot be after end time")