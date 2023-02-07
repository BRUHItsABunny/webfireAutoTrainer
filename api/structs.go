package api

import "strconv"

type WFTClass struct {
	Name      string
	Referer   string
	SessionID string
	Completed bool
}

type WFTCourse struct {
	URL     string
	Name    string
	Classes map[string]*WFTClass
}

type AICCData struct {
	LessonStatus string `json:"Lesson_Status"`
	Score        string `json:"Score"`
}

func NewAICCData(status LessonStatus) *AICCData {
	return &AICCData{
		LessonStatus: status.ToParameterString(),
		Score:        "",
	}
}

func NewAICCDataWithScore(status LessonStatus, score float64) *AICCData {
	return &AICCData{
		LessonStatus: status.ToParameterString(),
		Score:        strconv.FormatFloat(score, 'f', 2, 64),
	}
}
