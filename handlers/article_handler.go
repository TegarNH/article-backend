package handlers

import (
    "database/sql"
	"fmt"
    "net/http"
    "strconv"
    "time"
    "log"

    "github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
    "github.com/TegarNH/article-backend/models"
)

func formatValidationErrors(err error) []string {
	var errorMessages []string

	for _, e := range err.(validator.ValidationErrors) {
		// Menggunakan e.Tag() untuk mengetahui jenis validasi yang gagal
		var message string
		switch e.Tag() {
		case "required":
			message = fmt.Sprintf("%s tidak boleh kosong", e.Field())
		case "min":
			message = fmt.Sprintf("%s minimal %s karakter", e.Field(), e.Param())
		case "oneof":
			message = fmt.Sprintf("%s harus salah satu dari: %s", e.Field(), e.Param())
		default:
			message = fmt.Sprintf("Error pada field %s dengan tag %s", e.Field(), e.Tag())
		}
		errorMessages = append(errorMessages, message)
	}
	return errorMessages
}

type ArticleHandler struct {
    DB *sql.DB
}

func NewArticleHandler(db *sql.DB) *ArticleHandler {
    return &ArticleHandler{DB: db}
}

// 1. Create Article
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
    var newArticle models.Article

	if err := c.ShouldBindJSON(&newArticle); err != nil {
        // Cek apakah error ini adalah error validasi
        if verrors, ok := err.(validator.ValidationErrors); ok {
            errors := formatValidationErrors(verrors)
            c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
            return
        }
		// Jika bukan error validasi, kirim pesan umum
        c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Format JSON tidak valid"}})
        return
    }

    query := "INSERT INTO posts (title, content, category, status) VALUES (?, ?, ?, ?)"
    result, err := h.DB.Exec(query, newArticle.Title, newArticle.Content, newArticle.Category, newArticle.Status)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{"Failed to create article"}})
		log.Println("Error creating article:", err)
        return
    }

    id, _ := result.LastInsertId()
    newArticle.ID = int(id)
    newArticle.CreatedDate = time.Now()
    newArticle.UpdatedDate = time.Now()

    c.JSON(http.StatusCreated, newArticle)
}

// 2. Get Articles with Pagination
func (h *ArticleHandler) GetArticles(c *gin.Context) {
    limitStr := c.DefaultQuery("limit", "10")
    offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid limit parameter"}})
        return
    }

    offset, err := strconv.Atoi(offsetStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Invalid offset parameter"}})
        return
    }

    query := "SELECT id, title, content, category, created_date, updated_date, status FROM posts WHERE status != 'trash' LIMIT ? OFFSET ?"
    rows, err := h.DB.Query(query, limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{"Failed to fetch articles"}})
		log.Println("Error fetch article:", err)
        return
    }
    defer rows.Close()

    articles := []models.Article{}
    for rows.Next() {
        var article models.Article
        if err := rows.Scan(&article.ID, &article.Title, &article.Content, &article.Category, &article.CreatedDate, &article.UpdatedDate, &article.Status); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{"Failed to scan article data"}})
            return
        }
        articles = append(articles, article)
    }

    c.JSON(http.StatusOK, articles)
}

// 3. Get Article By ID
func (h *ArticleHandler) GetArticleByID(c *gin.Context) {
    id := c.Param("id")
    var article models.Article

    query := "SELECT id, title, content, category, created_date, updated_date, status FROM posts WHERE id = ? AND status != 'trash'"
    err := h.DB.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Content, &article.Category, &article.CreatedDate, &article.UpdatedDate, &article.Status)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"errors": []string{"Article not found"}})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{"Failed to fetch article"}})
		log.Println("Error fetch article:", err)
        return
    }

    c.JSON(http.StatusOK, article)
}

// 4. Update Article
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
    id := c.Param("id")
    var updatedArticle models.Article

    if err := c.ShouldBindJSON(&updatedArticle); err != nil {
        // Cek apakah error ini adalah error validasi
        if verrors, ok := err.(validator.ValidationErrors); ok {
            errors := formatValidationErrors(verrors)
            c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
            return
        }
		// Jika bukan error validasi, kirim pesan umum
        c.JSON(http.StatusBadRequest, gin.H{"errors": []string{"Format JSON tidak valid"}})
        return
    }

    query := "UPDATE posts SET title = ?, content = ?, category = ?, status = ? WHERE id = ?"
    _, err := h.DB.Exec(query, updatedArticle.Title, updatedArticle.Content, updatedArticle.Category, updatedArticle.Status, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{"Failed to update article"}})
		log.Println("Error update article:", err)
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Article " + id + " updated successfully"})
}

// 5. Delete Article
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
    id := c.Param("id")

    query := "UPDATE posts SET status = 'trash' WHERE id = ?"
    _, err := h.DB.Exec(query, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{"Failed to delete article"}})
		log.Println("Error delete article:", err)
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Article " + id + " deleted successfully"})
}