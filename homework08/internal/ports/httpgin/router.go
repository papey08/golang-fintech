package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework8/internal/app"
	"io"
)

func AppRouter(r *gin.RouterGroup, a app.App, logsWriter io.Writer) {
	r.Use(LoggerMiddleware(logsWriter))
	r.Use(PanicMiddleware())

	r.POST("/ads", createAd(a))                    // Метод для создания объявления (ad)
	r.PUT("/ads/:ad_id/status", changeAdStatus(a)) // Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
	r.PUT("/ads/:ad_id", updateAd(a))              // Метод для обновления текста(Text) или заголовка(Title) объявления
	r.GET("/ads/:ad_id", getAdByID(a))             // Метод для получения объявления по его ID
	r.GET("/ads", getAdsList(a))                   // Метод для получения списка объявлений с фильтрами
	r.GET("/ads/search", searchAds(a))             // Метод для поиска объявлений по названию

	r.GET("/users/:user_id", getUserByID(a)) // Метод для получения пользователя по его ID
	r.POST("/users", createUser(a))          // Метод для создания пользователя (user)
	r.PUT("/users/:user_id", updateUser(a))  // Метод для обновления nickname или email пользователя
}
