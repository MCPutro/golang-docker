package util

import (
	"errors"
	"log"
	"sync"
)

var (
	once           sync.Once
	ErrNotFound    error
	ErrAlreadyUsed error
	ErrNotMatch    error
)

func init() {
	once.Do(func() {
		ErrNotFound = errors.New("not found")
		ErrAlreadyUsed = errors.New("already used")
		ErrNotMatch = errors.New("do not match")
		log.Printf("init error type is OK")
	})
}
