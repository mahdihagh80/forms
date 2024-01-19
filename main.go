package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/mahdihagh80/forms/internal/controllers"
	"github.com/mahdihagh80/forms/internal/models"
	"github.com/mahdihagh80/forms/internal/services"
	"github.com/redis/go-redis/v9"
)

func main() {
	db, err := sqlx.Connect("mysql", "root:password@(localhost:3306)/test")
	if err != nil {
		log.Fatalln(err)
	}
	initializeDB(db)

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	store := models.NewUserModel(db)
	us := services.NewUserService(store)
	sStore := models.NewSessionStore(client)
	sService := services.NewSessionService(sStore)

	router := gin.Default()
	controllers.SetupUserController(us, sService)
	controllers.SetupRoutes(router)

	router.Run(":8080")

}

func initializeDB(db *sqlx.DB) {
	var content []byte
	var err error

	content, err = os.ReadFile("./migrations/users.sql")
	if err != nil {
		log.Fatal(err)
	}
	db.MustExec(string(content))

}

// fmt.Println(time.Now())
// fmt.Println(db.Exec("insert into sessions values (2, 1, ?)", time.Now()))

// var t string
// db.QueryRow("select expire from sessions where session_id = 2").Scan(&t)

// fmt.Println(t)
// expireTime, err := time.Parse("2006-01-02 15:04:05", t)
// fmt.Println("expireTime : ", expireTime.Local(), err)
