package bookshelf

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Book struct {
	gorm.Model
	Name      string
	Authors   string
	Libraries []Library `gorm:"many2many:book_libraries"`
}

func (b *Book) createTable(db *gorm.DB) {
	db.DropTableIfExists(b)
	if err := db.Debug().CreateTable(b).Error; err != nil {
		println(err)
	}

	db.Model(&Book{}).AddIndex("book_authors_ndx", "authors")
	db.Model(&Book{}).AddIndex("book_name_ndx", "name")

	db.Debug().Table("book_libraries").AddForeignKey("book_id", "books(id)", "CASCADE", "RESTRICT")
	db.Debug().Table("book_libraries").AddForeignKey("library_id", "libraries(id)", "CASCADE", "RESTRICT")
}

func CreateTables(db *gorm.DB) {
	(&Library{}).createTable(db)
	(&Book{}).createTable(db)
}
