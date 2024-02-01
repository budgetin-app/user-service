package database

import (
	"errors"
	"fmt"
)

func HandleErrorDB(err error) error {
	fmt.Printf("database error: %v", err)
	return errors.New("database error")
}
