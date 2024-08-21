package util

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		logrus.Fatal(err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		logrus.Fatal(err)
	}
	return db, mock
}
