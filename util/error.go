package util

import (
	"errors"
	"log"
	"sync"
)

var (
	once        sync.Once
	ErrNotFound error
)

func init() {
	once.Do(func() {
		ErrNotFound = errors.New("not found")
		log.Printf("init error type")
	})
}
