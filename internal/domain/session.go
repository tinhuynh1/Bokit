package domain

import "github.com/gorilla/websocket"

const (
	SessionStatusWaiting  = "waiting"
	SessionStatusActive   = "active"
	SessionStatusFinished = "finished"
)

type SessionCreate struct {
	Code   string `json:"code"`
	QuizId string `json:"quiz_id"`
	Status string `json:"status"`
}

type Session struct {
	Code              string                  `json:"code"`
	Quiz              *Quiz                   `json:"quiz"`
	Status            string                  `json:"status"`
	Participants      map[string]*Participant `json:"participants"`
	CurrentQuestion   int                     `json:"current_question"`
	StartTime         int64                   `json:"start_time"`
	EndTime           int64                   `json:"end_time"`
	QuestionStartedAt int64                   `json:"-"`
}
type Participant struct {
	ID       string          `json:"id"`
	Code     string          `json:"code"`
	Name     string          `json:"name"`
	Score    int             `json:"score"`
	Answers  map[int]string  `json:"answers"` // question index -> answer choice id
	Conn     *websocket.Conn `json:"-"`
	JoinedAt int64           `json:"joined_at"`
}

type LeaderboardEntry struct {
	ParticipantID string `json:"participant_id"`
	Name          string `json:"name"`
	Score         int    `json:"score"`
	Rank          int    `json:"rank"`
}
