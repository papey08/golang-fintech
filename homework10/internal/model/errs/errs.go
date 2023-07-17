package errs

import "errors"

var ValidationError = errors.New("invalid ad struct")
var WrongUserError = errors.New("different userID from ad.UserID")
var UserRepositoryError = errors.New("something wrong with users repository")
var AdRepositoryError = errors.New("something wrong with ads repository")
var AdNotExist = errors.New("no ad in repository with required id")
var UserNotExist = errors.New("no users in repository with required id")
