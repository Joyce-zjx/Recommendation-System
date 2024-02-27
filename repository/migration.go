package repository

import (
	"Recommendation-System/repository/schema"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type GormDatabase struct {
	DB *gorm.DB
}

func NewGormDatabase(dsn string, debug bool) (*GormDatabase, error) {
	config := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	if debug {
		config.Logger = gormLogger.Default.LogMode(gormLogger.Info)
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, err
	}
	return &GormDatabase{DB: db}, nil
}

func (d *GormDatabase) AutoMigrate() error {
	// enable format UUID as PK
	d.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err := d.DB.AutoMigrate(
		&schema.User{},
		&schema.UserLikeGenreAndSubGenre{},
		&schema.UserDisLikeGenreAndSubGenre{},
	); err != nil {
		return err
	}
	return nil
}
