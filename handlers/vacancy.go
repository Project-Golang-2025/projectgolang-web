package handlers

import (
    "net/http"
    "strings"
    "vacancy_api/models"
    "vacancy_api/storage"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

func GetVacancies(c *gin.Context) {
    query := strings.ToLower(c.Query("q"))
    vacancies := storage.GetVacancies()
    if query == "" {
        c.JSON(http.StatusOK, vacancies)
        return
    }
    filtered := []models.Vacancy{}
    for _, v := range vacancies {
        if strings.Contains(strings.ToLower(v.Title), query) ||
            strings.Contains(strings.ToLower(v.Company), query) ||
            strings.Contains(strings.ToLower(v.Description), query) ||
            strings.Contains(strings.ToLower(v.Status), query) ||
            strings.Contains(strings.ToLower(v.ExperienceLevel), query) {
            filtered = append(filtered, v)
            continue
        }
        for _, kw := range v.Keywords {
            if strings.Contains(strings.ToLower(kw), query) {
                filtered = append(filtered, v)
                break
            }
        }
    }
    c.JSON(http.StatusOK, filtered)
}

func AddVacancy(c *gin.Context) {
    var v models.Vacancy
    if err := c.ShouldBindJSON(&v); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    v.ID = uuid.New().String()
    newVac, err := storage.AddVacancy(v)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, newVac)
}

func UpdateVacancy(c *gin.Context) {
    id := c.Param("id")
    var v models.Vacancy
    if err := c.ShouldBindJSON(&v); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    updated, err := storage.UpdateVacancyByID(id, v)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Vacancy not found"})
        return
    }
    c.JSON(http.StatusOK, updated)
}

func DeleteVacancy(c *gin.Context) {
    id := c.Param("id")
    err := storage.DeleteVacancyByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Vacancy not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"result": "deleted"})
}
