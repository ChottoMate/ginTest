package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	gorm.Model
	Text   string
	Status string
}

func dbInit() {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("database don't open. (dbInit)")
	}
	db.AutoMigrate(&Todo{})
	defer db.Close()
}

func dbInsert(text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("database don't open. (dbInsert)")
	}
	db.Create(&Todo{Text: text, Status: status})
	defer db.Close()
}

func dbUpdate(id int, text string, status string) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("database don't open. (dbUpdate)")
	}
	var todo Todo
	db.First(&todo, id)
	todo.Text = text
	todo.Status = status
	db.Save(&todo)
	db.Close()
}

func dbDelete(id int) {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("database don't open. (dbDelete)")
	}
	var todo Todo
	db.First(&todo, id)
	db.Delete(&todo)
	db.Close()
}

func dbFindAll() []Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("database don't open")
	}
	var todos []Todo
	db.Order("created_at desc").Find(&todos)
	db.Close()
	return todos
}

func dbGetOne(id int) Todo {
	db, err := gorm.Open("sqlite3", "test.sqlite3")
	if err != nil {
		panic("database don't open. (dbGetOne)")
	}
	var todo Todo
	db.First(&todo, id)
	db.Close()
	return todo
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	dbInit()

	//index
	router.GET("/", func(context *gin.Context) {
		todos := dbFindAll()
		context.HTML(200, "index.html", gin.H{"todos": todos})
	})

	//create
	router.POST("/new", func(context *gin.Context) {
		text := context.PostForm("text")
		status := context.PostForm("status")
		dbInsert(text, status)
		context.Redirect(302, "/")
	})

	//detail
	router.GET("/detail/:id", func(context *gin.Context) {
		n := context.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		todo := dbGetOne(id)
		context.HTML(200, "detail.html", gin.H{"todo": todo})
	})

	//update
	router.POST("/update/:id", func(context *gin.Context) {
		n := context.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		text := context.PostForm("text")
		status := context.PostForm("status")
		dbUpdate(id, text, status)
		context.Redirect(302, "/")
	})

	//verify delete
	router.GET("/delete_check/:id", func(context *gin.Context) {
		n := context.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		todo := dbGetOne(id)
		context.HTML(200, "delete.html", gin.H{"todo": todo})
	})

	//delete
	router.POST("/delete/:id", func(context *gin.Context) {
		n := context.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		dbDelete(id)
		context.Redirect(302, "/")
	})

	router.Run()
}
