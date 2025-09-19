package domain

import "time"

// Event Types
const (
	EventTypeLeaderboardUpdated = "leaderboard.updated"
	EventTypeParticipantJoined  = "participant.joined"
	EventTypeParticipantLeft    = "participant.left"
	EventTypeSessionCreated     = "session.created"
	EventTypeSessionStarted     = "session.started"
	EventTypeSessionEnded       = "session.ended"
	EventTypeAnswerSubmitted    = "answer.submitted"
)

// Base Event Structure
type BaseEvent struct {
	Type      string    `json:"type"`
	SessionID string    `json:"session_id"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}

// Leaderboard Updated Event
type LeaderboardUpdatedEvent struct {
	BaseEvent
	Data LeaderboardUpdatedData `json:"data"`
}

type LeaderboardUpdatedData struct {
	ParticipantID string `json:"participant_id"`
	ParticipantName string `json:"participant_name"`
	Score         int    `json:"score"`
	Rank          int    `json:"rank"`
	Leaderboard   []LeaderboardEntry `json:"leaderboard"`
}

// Participant Joined Event
type ParticipantJoinedEvent struct {
	BaseEvent
	Data ParticipantJoinedData `json:"data"`
}

type ParticipantJoinedData struct {
	ParticipantID   string `json:"participant_id"`
	ParticipantName string `json:"participant_name"`
	TotalParticipants int  `json:"total_participants"`
}

// Participant Left Event
type ParticipantLeftEvent struct {
	BaseEvent
	Data ParticipantLeftData `json:"data"`
}

type ParticipantLeftData struct {
	ParticipantID   string `json:"participant_id"`
	ParticipantName string `json:"participant_name"`
	TotalParticipants int  `json:"total_participants"`
}

// Session Created Event
type SessionCreatedEvent struct {
	BaseEvent
	Data SessionCreatedData `json:"data"`
}

type SessionCreatedData struct {
	Code   string `json:"code"`
	QuizID string `json:"quiz_id"`
	Quiz   *Quiz  `json:"quiz"`
}

// Session Started Event
type SessionStartedEvent struct {
	BaseEvent
	Data SessionStartedData `json:"data"`
}

type SessionStartedData struct {
	TotalParticipants int `json:"total_participants"`
	Quiz             *Quiz `json:"quiz"`
}

// Session Ended Event
type SessionEndedEvent struct {
	BaseEvent
	Data SessionEndedData `json:"data"`
}

type SessionEndedData struct {
	FinalLeaderboard []LeaderboardEntry `json:"final_leaderboard"`
	TotalParticipants int `json:"total_participants"`
}

// Answer Submitted Event
type AnswerSubmittedEvent struct {
	BaseEvent
	Data AnswerSubmittedData `json:"data"`
}

type AnswerSubmittedData struct {
	ParticipantID   string `json:"participant_id"`
	ParticipantName string `json:"participant_name"`
	QuestionIndex   int    `json:"question_index"`
	AnswerChoice    string `json:"answer_choice"`
	Correct         bool   `json:"correct"`
	Score           int    `json:"score"`
	NewTotalScore   int    `json:"new_total_score"`
}
