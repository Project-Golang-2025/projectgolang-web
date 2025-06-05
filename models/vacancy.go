package models

type Vacancy struct {
    ID              string   `json:"id"`
    Title           string   `json:"title"`
    Company         string   `json:"company"`
    Description     string   `json:"description"`
    Keywords        []string `json:"keywords"`
    SourceURL       string   `json:"sourceURL,omitempty"`
    Status          string   `json:"status,omitempty"`
    ExperienceLevel string   `json:"experienceLevel,omitempty"`
    Notes           string   `json:"notes,omitempty"`
    ResumePath      string   `json:"resumePath,omitempty"`
    ResumeFileName  string   `json:"resumeFileName,omitempty"`
}
