package database

import "gorm.io/gorm"

// SeederDB seeds the database with the defined values
// Note: Use 'Save' method when seeding the db, so it doesn't create duplicates
func SeederDB(db *gorm.DB) {
	// .. add db seeder here
}
