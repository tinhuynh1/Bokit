package dto

type CreateQuizRequest struct {
	Name      string         `json:"name"`
	Questions []QuizQuestion `json:"questions"`
}

type QuizQuestion struct {
	Name    string       `json:"name"`
	Time    int          `json:"time"`
	Choices []QuizChoice `json:"choices"`
}

type QuizChoice struct {
	Name    string `json:"name"`
	Correct bool   `json:"correct"`
}
