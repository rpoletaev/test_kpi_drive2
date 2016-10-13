package bookshelf

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Library struct {
	gorm.Model
	Name      string
	Libraries []Library `gorm:"many2many:book_libraries"`
}

func (l *Library) createTable(db *gorm.DB) {
	db.DropTableIfExists(l)
	db.CreateTable(l)
}
