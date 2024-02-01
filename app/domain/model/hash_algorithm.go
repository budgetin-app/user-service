package model

type HashAlgorithm struct {
	ID   uint   `gorm:"column:hash_algorithm_id; primaryKey"`
	Name string `gorm:"column:algorithm_name; size:20; unique"`
	BaseModel
}
