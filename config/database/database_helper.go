package database

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func HandleErrorDB(err error) error {
	log.Errorf("database error: %v", err)
	return fmt.Errorf("database error: %v", err)
}
