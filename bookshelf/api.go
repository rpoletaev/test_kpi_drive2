package bookshelf

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/tommy351/gin-csrf"
	"github.com/tommy351/gin-sessions"
	"github.com/unrolled/render"
)

var rnd *render.Render

type Config struct {
	Port                string `yaml:"port"`
	PostgeressConString string `yaml:"connection_string"`
	Secret              string `yaml:"secret"`
}

type API struct {
	router *gin.Engine
	config *Config
	log    *log.Logger
	db     *gorm.DB
}

func (api *API) setRoutes() {
	api.router.GET("/", api.bookshelf)
	api.router.GET("/book", api.books)
	api.router.GET("/book/:id", api.getBook)
	api.router.POST("/book", api.createBook)
	api.router.PUT("/book", api.editBook)
	api.router.DELETE("book", api.deleteBook)

	api.router.GET("/libs", api.libs)
	api.router.GET("/lib/:id", api.getLib)
	api.router.POST("/lib", api.createLib)
	api.router.PUT("/lib", api.editLib)
	api.router.DELETE("lib", api.deleteLib)
}

//NewAPI Create & configure API object
func NewAPI(config Config, init bool) (api *API, err error) {
	api = &API{
		config: &config,
		router: gin.Default(),
		log:    log.New(os.Stdout, "api", -1),
	}
	var db *gorm.DB
	store := sessions.NewCookieStore([]byte(api.config.Secret))
	db, err = gorm.Open("postgres", config.PostgeressConString)
	if err != nil {
		return nil, err
	}
	api.db = db

	if init {
		api.initialization()
	}

	api.router.Use(gin.Logger())
	api.router.Use(gin.Recovery())
	api.router.Use(func(c *gin.Context) {
		println("FROM MY NEW MIDDLEWARE")
		println(len(c.Request.PostForm))
		for k, v := range c.Request.PostForm {
			fmt.Println("key:", k)
			fmt.Println("val:", strings.Join(v, ""))
		}
	})
	api.router.Use(sessions.Middleware("my_sessions", store))
	api.router.Use(csrf.Middleware(csrf.Options{Secret: api.config.Secret, IgnoreMethods: []string{"GET"}}))
	api.router.Static("/static", "./static")
	api.router.LoadHTMLGlob("views/*")
	api.setRoutes()

	return api, nil
}

func (api *API) bookshelf(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"_csrf": csrf.GetToken(c)})
}

func (api *API) initialization() {
	api.log.Println("Iint tables")
	CreateTables(api.db)
}

func (api *API) Run() {
	log.Fatal(api.router.Run(api.config.Port))
}

func PostForm(c *gin.Context, key string) string {
	return template.HTMLEscapeString(c.PostForm(key))
}

func GetPostForm(c *gin.Context, key string) (string, bool) {
	val, exist := c.GetPostForm(key)
	if !exist {
		return val, exist
	}

	val = template.HTMLEscapeString(val)
	return val, exist
}
