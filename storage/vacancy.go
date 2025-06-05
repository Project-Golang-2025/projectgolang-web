package storage

import (
    "bufio"
    "encoding/json"
    "os"
    "sync"
    "vacancy_api/models"
    "log"
    "time"
)

const VacanciesFile = "vacancies.jsonl"

var (
    vacancies      []models.Vacancy
    vacanciesMutex sync.RWMutex

    saveChan = make(chan struct{}, 1)
)

func init() {
    go saveWorker()
}

func saveWorker() {
    for range saveChan {
        err := saveAllVacanciesInternal()
        if err != nil {
            log.Printf("Ошибка сохранения вакансий: %v", err)
        }
        time.Sleep(100 * time.Millisecond)
    }
}

func LoadVacancies() error {
    vacanciesMutex.Lock()
    defer vacanciesMutex.Unlock()

    file, err := os.Open(VacanciesFile)
    if err != nil {
        if os.IsNotExist(err) {
            vacancies = []models.Vacancy{}
            return nil
        }
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    var loaded []models.Vacancy
    for scanner.Scan() {
        var v models.Vacancy
        if err := json.Unmarshal(scanner.Bytes(), &v); err != nil {
            continue
        }
        loaded = append(loaded, v)
    }
    vacancies = loaded
    return scanner.Err()
}

func GetVacancies() []models.Vacancy {
    vacanciesMutex.RLock()
    defer vacanciesMutex.RUnlock()
    return append([]models.Vacancy(nil), vacancies...)
}

func AddVacancy(v models.Vacancy) (models.Vacancy, error) {
    vacanciesMutex.Lock()
    vacancies = append(vacancies, v)
    vacanciesMutex.Unlock()

    file, err := os.OpenFile(VacanciesFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return v, err
    }
    defer file.Close()

    data, err := json.Marshal(v)
    if err != nil {
        return v, err
    }
    data = append(data, '\n')
    if _, err := file.Write(data); err != nil {
        return v, err
    }

    return v, nil
}

func saveAllVacanciesInternal() error {
    vacanciesMutex.RLock()
    defer vacanciesMutex.RUnlock()

    file, err := os.Create(VacanciesFile)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    for _, v := range vacancies {
        data, err := json.Marshal(v)
        if err != nil {
            continue
        }
        data = append(data, '\n')
        if _, err := writer.Write(data); err != nil {
            return err
        }
    }
    return writer.Flush()
}

func requestSave() {
    select {
    case saveChan <- struct{}{}:
    default:
    }
}

func UpdateVacancyByID(id string, v models.Vacancy) (models.Vacancy, error) {
    vacanciesMutex.Lock()
    defer vacanciesMutex.Unlock()
    for i := range vacancies {
        if vacancies[i].ID == id {
            v.ID = id
            vacancies[i] = v
            requestSave()
            return v, nil
        }
    }
    return v, os.ErrNotExist
}

func DeleteVacancyByID(id string) error {
    vacanciesMutex.Lock()
    defer vacanciesMutex.Unlock()
    for i := range vacancies {
        if vacancies[i].ID == id {
            vacancies = append(vacancies[:i], vacancies[i+1:]...)
            requestSave()
            return nil
        }
    }
    return os.ErrNotExist
}
