package domain

// WebSocket Packet Types
const (
	PacketJoinSession    = 1
	PacketLeaveSession   = 2
	PacketAnswerQuestion = 3
	PacketStartQuiz      = 4
	PacketNextQuestion   = 5
	PacketQuizEnded      = 6
	PacketLeaderboard    = 7
	PacketSessionUpdate  = 8
	PacketError          = 9
)

// Join Session Request
type JoinSessionRequest struct {
	SessionCode     string `json:"session_code"`
	ParticipantName string `json:"participant_name"`
}

// Join Session Response
type JoinSessionResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Session *Session `json:"session,omitempty"`
}

// Answer Question Request
type AnswerQuestionRequest struct {
	QuestionIndex int    `json:"question_index"`
	AnswerChoice  string `json:"answer_choice"`
}

// Answer Question Response
type AnswerQuestionResponse struct {
	Success       bool `json:"success"`
	Correct       bool `json:"correct"`
	Score         int  `json:"score"`
	QuestionIndex int  `json:"question_index"`
}

// Start Quiz Response
type StartQuizResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Quiz    *Quiz  `json:"quiz,omitempty"`
}

// Next Question Response
type NextQuestionResponse struct {
	QuestionIndex int           `json:"question_index"`
	Question      *QuizQuestion `json:"question"`
	TimeLeft      int           `json:"time_left"`
}

// Quiz Ended Response
type QuizEndedResponse struct {
	FinalLeaderboard []LeaderboardEntry `json:"final_leaderboard"`
	YourScore        int                `json:"your_score"`
	YourRank         int                `json:"your_rank"`
}

// Leaderboard Update
type LeaderboardUpdate struct {
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
}

// Session Update
type SessionUpdate struct {
	Session *Session `json:"session"`
}

// Error Response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
