# JSON-CRUD-API-in-Go
JSON CRUD API in Go (Gin/GORM)

ref: [Creating a JSON CRUD API in Go (Gin/GORM)](https://www.youtube.com/watch?v=lf_kiH_NPvM&ab_channel=CodingwithRobby)

- package:
    - [Gin](https://gin-gonic.com/)
    - [Gorm](https://gorm.io/)
    - [githubnemo/CompileDaemon](https://github.com/githubnemo/CompileDaemon)
    - [joho/godotenv](https://github.com/joho/godotenv)
    
- software(free)
    - [ElephantSQL](https://customer.elephantsql.com/)
    - [Postman](https://www.postman.com/)
    - [TablePlus](https://tableplus.com/)

---

## Declaring Models
- Models are normal structs with basic Go types, pointers/alias of them or custom types implementing Scanner and Valuer interfaces
```go
// postModel.go
package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title string
	Body  string
}

// gorm.Model definition
// type Model struct {
// 	ID        uint           `gorm:"primaryKey"`
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// 	DeletedAt gorm.DeletedAt `gorm:"index"`
//   }
```
## Database

### Create Database
![](https://i.imgur.com/q1eb8b0.png)
#### Get the DB URL
![](https://i.imgur.com/tld974g.png)
- Paste it in the `.env`
    - can see in `9:59` on video
- Default format

```go
DB_URL="host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
```


### Connecting to database
```go
// database.go
package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


var DB *gorm.DB

func ConnectToDB() {
	var err error

	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database.")
	}
}
```

- `os.Getenv("DB_URL")` fetch `DB_URL` from `.env`

## Load env variables
```go
// loadEnvVar.go
package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEndVar() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
```

## Auto Migration
- Automatically migrate your schema, to keep your schema up to date.

```go
//migrate.go
package main

import (
	"go-crud/initializers"
	"go-crud/models"
)

// func init() will execute before func main()
func init() {
	initializers.LoadEndVar()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.Post{})
}

```
* Execute this file `go run migrate.go`
* Open `TablePlus`
    * ![](https://i.imgur.com/1pVZpVM.png)
    * Paste the URL from elephantSQL
    * click `import` and `connect`
    * ![](https://i.imgur.com/rJwKlld.png)
    * Create table successfully


### CRUD

### Create data
```go
//postController.go
package controllers

import (
	"go-crud/initializers"
	"go-crud/models"

	"github.com/gin-gonic/gin"
)

func PostsCreate(c *gin.Context) {
	// Get data off req body
	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	// Create a post
	post := models.Post{Title: body.Title, Body: body.Body}

	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.Status(400)
		return
	}
	// Return it
	c.JSON(200, gin.H{
		"post": post,
	})
}
```
```go
// main.go
func main(){
    ...
    r.POST("/posts", controllers.PostsCreate)
    ...
}
```
#### Test on Postman and TablePlus
![](https://i.imgur.com/S4DR3fB.png)
![](https://i.imgur.com/FRsIMpJ.png)

---

### Retrieving all objects
```go
//postController.go
func PostsIndex(c *gin.Context) {
	// Get the posts
	var posts []models.Post
	initializers.DB.Find(&posts)

	// Respond with them
	c.JSON(200, gin.H{
		"posts": posts,
	})

}
```

```go
// main.go
func main(){
    ...
    r.GET("/posts", controllers.PostsIndex)
    ...
}
```

#### Test on Postman
![](https://i.imgur.com/1jDTi9L.png)

---

### Retrieving objects with ID
```go
//postController.go
func PostsShow(c *gin.Context) {
	// Get id off URL
	id := c.Param("id")

	// Get the posts
	var post models.Post
	initializers.DB.First(&post, id)

	// Respond with them
	c.JSON(200, gin.H{
		"post": post,
	})

}
```
```go
// main.go
func main() {
	...
	r.GET("/posts/:id", controllers.PostsShow)
    ...
}
```

#### Test on Postman
![](https://i.imgur.com/YyxfDdY.png)

---

### Update
```go
// postController.go
func PostsUpdate(c *gin.Context) {
	// Get id off URL
	id := c.Param("id")

	// Get the data off req body
	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	// Find the post were updating
	var post models.Post
	initializers.DB.First(&post, id)

	// Update it
	initializers.DB.Model(&post).Updates(models.Post{
		Title: body.Title,
		Body:  body.Body,
	})
	// Respond with them
	c.JSON(200, gin.H{
		"post": post,
	})

}
```
```go
// main.go
func main() {
    ...
    r.PUT("/posts/:id", controllers.PostsUpdate)
    ...
}
```

#### Test on Postman
- Previous
![](https://i.imgur.com/YyxfDdY.png)
- After
![](https://i.imgur.com/obO9m1a.png)

---

### Delete
```go
// postController.go
func PostsDelete(c *gin.Context) {
	// Get the id off URL
	id := c.Param("id")

	// Delete the post
	initializers.DB.Delete(&models.Post{}, id)

	// Repond
	c.Status(200)

}
```
```go
// main.go
func main() {
    ...
    r.DELETE("/posts/:id", controllers.PostsDelete)
    ...
}
```
#### Test om Postman and TablePlus
![](https://i.imgur.com/oFArOny.png)
![](https://i.imgur.com/rqwYUhF.png)


