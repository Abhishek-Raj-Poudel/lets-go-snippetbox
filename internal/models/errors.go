package models

import (
  "errors"
)

var ( 
  ErrNoRecord = errors.New("models: No Matching Records Found")

  ErrInvalidCredentials = errors.New("models: invalid credentials")

  ErrDuplicateEmail = errors.New("models: invalid credentials")
)
