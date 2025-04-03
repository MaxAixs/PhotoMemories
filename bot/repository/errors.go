package repository

import "errors"

var (
	PicNotFound  = errors.New("pictures not found in your collection for tag")
	PicExists    = errors.New("an image for the given tag already exists, please try again")
	UserNotFound = errors.New("you don't have any saved tags.⚠️")
)
