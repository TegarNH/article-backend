package main

import (
    "database/sql"
    "log"
	"os"

	"github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"  
    _ "github.com/go-sql-driver/mysql"
    "github.com/TegarNH/article-backend/handlers"
)

func main() {
	gin.SetMode(os.Getenv("GIN_MODE"))

	// load environment
	godotenv.Load()

    // Koneksi ke database MySQL
    dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default
	}

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal("Failed to ping database:", err)
    }

    log.Println("Database connected successfully")

    // Inisialisasi Gin router
    router := gin.Default()

	// Setup CORS
	router.Use(cors.Default())

    h := handlers.NewArticleHandler(db)

    // Grouping route untuk /article
    articleRoutes := router.Group("/article")
    {
        articleRoutes.POST("/", h.CreateArticle)                 // No 1: Membuat article baru
        articleRoutes.GET("/", h.GetArticles)      			     // No 2: Menampilkan seluruh article

        articleRoutes.GET("/:id", h.GetArticleByID)              // No 3: Menampilkan article by ID
        articleRoutes.PUT("/:id", h.UpdateArticle)               // No 4: Merubah data article
        articleRoutes.DELETE("/:id", h.DeleteArticle)            // No 5: Menghapus data article
    }

	router.Run(":" + port)
}