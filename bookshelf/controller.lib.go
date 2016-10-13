package bookshelf

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	csrf "github.com/tommy351/gin-csrf"
)

func (api *API) libs(c *gin.Context) {
	libs := []Library{}
	if err := api.db.Model(&Library{}).Find(&libs).Error; err != nil {
		api.log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список библиотек"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"libraries": libs, "_csrf": csrf.GetToken(c)})
}

func (api *API) getLib(c *gin.Context) {

}

func (api *API) createLib(c *gin.Context) {
	if name, ok := GetPostForm(c, "name"); ok {
		newLib := &Library{Name: name}
		if err := api.db.Create(newLib).Error; err != nil {
			api.log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать библиотеку", "_csrf": csrf.GetToken(c)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"_csrf": csrf.GetToken(c)})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Необходимо указать имя новой библиотеки", "_csrf": csrf.GetToken(c)})
}

func (api *API) editLib(c *gin.Context) {
	lib := &Library{}
	val, exist := GetPostForm(c, "id")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Необходимо указать код изменяемой библиотеки", "_csrf": csrf.GetToken(c)})
		return
	}

	id, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректно указан код библиотеки", "_csrf": csrf.GetToken(c)})
		return
	}

	lib.ID = uint(id)

	val, exist = GetPostForm(c, "name")
	if !exist || strings.TrimSpace(val) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Необходимо указать новое имя библиотеки", "_csrf": csrf.GetToken(c)})
		return
	}

	if err = api.db.Model(lib).Update("name", val).Error; err != nil {
		api.log.Println("Library update error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить имя библиотеки", "_csrf": csrf.GetToken(c)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"_csrf": csrf.GetToken(c)})
}

func (api *API) deleteLib(c *gin.Context) {
	val, exist := GetPostForm(c, "id")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан код библиотеки", "_csrf": csrf.GetToken(c)})
		return
	}

	id, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректно указан код библиотеки", "_csrf": csrf.GetToken(c)})
		return
	}

	if err := api.db.Debug().Delete(&Library{}, "id = ?", uint(id)).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось удалить библиотеку", "_csrf": csrf.GetToken(c)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"_csrf": csrf.GetToken(c)})
}
