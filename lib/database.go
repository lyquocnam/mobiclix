package lib

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go-mobiclix/app/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	DB, err = gorm.Open("postgres", `host=localhost port=5433 dbname=mobiclix user=postgres password=postgress sslmode=disable`)
	//DB, err = gorm.Open("mysql", "root:root@/mobiclix?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic(err)
	}
	// DB.LogMode(true)
	DB.DB().SetMaxIdleConns(0)
	migrate()
	seed()
}

func migrate() {
	if err := DB.AutoMigrate(&models.Ticket{}).Error; err != nil {
		panic(err)
	}
}

func seed() {
	count := 0
	if err := DB.Model(&models.Ticket{}).Count(&count).Error; err != nil {
		panic(err)
	}

	if count > 0 {
		return
	}

	tx := DB.Begin()
	for i := 0; i < 1000; i++ {
		if err := tx.Create(&models.Ticket{}).Error; err != nil {
			panic(err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		panic(err)
	}
}
