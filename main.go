package main

import (
    "log"
    "time"
    "vacancy_api/handlers"
    "vacancy_api/storage"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    if err := storage.LoadVacancies(); err != nil {
        log.Fatalf("Ошибка загрузки вакансий: %v", err)
    }

    r := gin.Default()

    // Настройка CORS
    config := cors.Config{
        AllowAllOrigins:  true, // В продакшене лучше указать конкретные домены
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }
    r.Use(cors.New(config))

    v := r.Group("/api/vacancies")
    {
        v.GET("", handlers.GetVacancies)
        v.POST("", handlers.AddVacancy)
        v.PUT("/:id", handlers.UpdateVacancy)
        v.DELETE("/:id", handlers.DeleteVacancy)
    }

    log.Println("Сервер запущен на :8081")
    r.Run(":8081")
}
