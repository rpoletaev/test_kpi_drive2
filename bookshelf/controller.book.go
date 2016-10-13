package bookshelf

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rpoletaev/test_kpi_drive2/bookshelf/utils"
	"github.com/tommy351/gin-csrf"
)

func (api *API) books(c *gin.Context) {
	books := []Book{}

	q := c.Request.URL.Query()
	libVal := q.Get("lib")
	lib, err := strconv.ParseInt(libVal, 10, 32)

	if err == nil && lib > 0 {
		api.db.Debug().Joins("JOIN book_libraries on book_libraries.book_id = books.id AND book_libraries.library_id = ?", lib).Find(&books)
		c.JSON(http.StatusOK, gin.H{"books": books, "_csrf": csrf.GetToken(c)})
		return
	}

	if err := api.db.Find(&books).Error; err != nil {
		api.log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"status": "Не удалось получить список книг", "_csrf": csrf.GetToken(c)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"books": books, "_csrf": csrf.GetToken(c)})
}

func (api *API) getBook(c *gin.Context) {
	idVal := c.Param("id")
	id, err := strconv.ParseInt(idVal, 10, 32)
	if err != nil {
		api.log.Println(err.Error)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Неверно указан код книги", "_csrf": csrf.GetToken(c)})
		return
	}

	book := Book{}
	if dberr := api.db.Find(&book, id).Error; dberr != nil {
		api.log.Println(dberr.Error)
		c.JSON(http.StatusNotFound, gin.H{"status": "Не удалось найти книгу с указанным кодом", "_csrf": csrf.GetToken(c)})
		return
	}

	libraries := []uint{}
	api.db.Table("book_libraries").Where("book_id = ?", id).Pluck("library_id", &libraries)
	api.log.Printf("%v\r\n", libraries)
	c.JSON(http.StatusOK, gin.H{"book": book, "libraries": libraries, "_csrf": csrf.GetToken(c)})
	return
}

func (api *API) createBook(c *gin.Context) {
	name := PostForm(c, "name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Должно быть указано название книги", "_csrf": csrf.GetToken(c)})
		return
	}

	authors := PostForm(c, "authors")
	if authors == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Должен быть указан автор", "_csrf": csrf.GetToken(c)})
		return
	}

	val, exist := GetPostForm(c, "libraries")
	if !exist || strings.TrimSpace(val) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Не выбрано ни одной библиотеки", "_csrf": csrf.GetToken(c)})
		return
	}

	libIds := utils.StringToSliceUI(val)
	libraries := []Library{}
	if err := api.db.Find(&libraries, "id in (?)", libIds).Error; err != nil {
		api.log.Println("Ошибка при получении списка библиотек\n", err.Error)
	}

	if err := api.db.Model(&Book{}).Create(&Book{
		Name:      name,
		Authors:   authors,
		Libraries: libraries,
	}).Error; err != nil {
		api.log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Не удалось создать книгу", "_csrf": csrf.GetToken(c)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"_csrf": csrf.GetToken(c)})
}

func (api *API) editBook(c *gin.Context) {
	book := Book{}
	val, exist := GetPostForm(c, "id")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Не указан код книги", "_csrf": csrf.GetToken(c)})
		return
	}

	id, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Некорректно указан код книги", "_csrf": csrf.GetToken(c)})
		return
	}

	book.ID = uint(id)

	val, exist = GetPostForm(c, "name")
	if !exist || strings.TrimSpace(val) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Не указано название книги", "_csrf": csrf.GetToken(c)})
		return
	}
	book.Name = val

	val, exist = GetPostForm(c, "authors")
	if !exist || strings.TrimSpace(val) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Не указан автор книги", "_csrf": csrf.GetToken(c)})
		return
	}
	book.Authors = val

	val, exist = GetPostForm(c, "libraries")
	if !exist || strings.TrimSpace(val) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Не выбрано ни одной библиотеки", "_csrf": csrf.GetToken(c)})
		return
	}

	libIds := utils.StringToSliceUI(val)
	libraries := []Library{}
	if err = api.db.Find(&libraries, "id in (?)", libIds).Error; err != nil {
		api.log.Println("Ошибка при получении списка библиотек\n", err.Error)
	}

	book.Libraries = libraries
	if err = api.db.Save(&book).Error; err != nil {
		api.log.Println("Book update error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Не удалось обновить имя библиотеки", "_csrf": csrf.GetToken(c)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"_csrf": csrf.GetToken(c)})

}

func (api *API) deleteBook(c *gin.Context) {
	valId, ok := c.GetPostForm("id")
	if !ok {
		api.log.Println("не указан id для удаления книги")
		c.JSON(http.StatusBadRequest, gin.H{"status": "Для удаления необходимо указать код", "_csrf": csrf.GetToken(c)})
		return
	}

	intID, err := strconv.Atoi(valId)
	if err != nil {
		api.log.Printf("Указан некорректный код: %s\n", valId)
		c.JSON(http.StatusBadRequest, gin.H{"status": "Указан некорректный код", "_csrf": csrf.GetToken(c)})
		return
	}

	if err = api.db.Delete(&Book{}, uint(intID)).Error; err != nil {
		api.log.Println("Ошибка при удалении книги с кодом: ", valId, "\n", err.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Не удалось удалить книгу", "_csrf": csrf.GetToken(c)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"_csrf": csrf.GetToken(c)})
}
