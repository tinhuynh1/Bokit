package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"quiz-svc/internal/domain"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type SessionService struct {
	sessions map[string]*domain.Session
	mutex    sync.RWMutex
}

func NewSessionService() *SessionService {
	return &SessionService{
		sessions: make(map[string]*domain.Session),
	}
}

func (s *SessionService) CreateSession(ctx context.Context, quiz *domain.Quiz) (*domain.Session, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session := &domain.Session{
		Code:            s.generateCode(),
		Quiz:            quiz,
		Status:          domain.SessionStatusWaiting,
		Participants:    make(map[string]*domain.Participant),
		CurrentQuestion: -1,
		StartTime:       0,
		EndTime:         0,
	}

	// Ensure unique code
	for {
		if _, exists := s.sessions[session.Code]; !exists {
			break
		}
		session.Code = s.generateCode()
	}

	s.sessions[session.Code] = session
	fmt.Println("Sessions", s.sessions)
	return session, nil
}

func (s *SessionService) generateCode() string {
	return strconv.Itoa(100000 + rand.Intn(900000))
}

func (s *SessionService) GetSession(ctx context.Context, code string) (*domain.Session, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if session, ok := s.sessions[code]; ok {
		return session, nil
	}
	return nil, errors.New("session not found")
}

func (s *SessionService) JoinSession(ctx context.Context, code, participantName string, conn *websocket.Conn) (*domain.Participant, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, exists := s.sessions[code]
	if !exists {
		return nil, errors.New("session not found")
	}

	if session.Status != domain.SessionStatusWaiting {
		return nil, errors.New("session is not accepting new participants")
	}

	participantID := uuid.New().String()
	participant := &domain.Participant{
		ID:       participantID,
		Name:     participantName,
		Score:    0,
		Answers:  make(map[int]string),
		Conn:     conn,
		JoinedAt: time.Now().Unix(),
	}

	session.Participants[participantID] = participant
	return participant, nil
}

func (s *SessionService) LeaveSession(ctx context.Context, code, participantID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, exists := s.sessions[code]
	if !exists {
		return errors.New("session not found")
	}

	delete(session.Participants, participantID)
	return nil
}

func (s *SessionService) StartSession(ctx context.Context, code string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	fmt.Println("Starting session", code)
	fmt.Println("Sessions", s.sessions)
	session, exists := s.sessions[code]
	if !exists {
		return errors.New("session not found")
	}

	if session.Status != domain.SessionStatusWaiting {
		return errors.New("session cannot be started")
	}

	if len(session.Participants) == 0 {
		return errors.New("no participants in session")
	}

	session.Status = domain.SessionStatusActive
	session.CurrentQuestion = 0
	session.StartTime = time.Now().Unix()
	return nil
}

func (s *SessionService) SubmitAnswer(ctx context.Context, code, participantID string, questionIndex int, answerChoice string) (*domain.AnswerQuestionResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, exists := s.sessions[code]
	if !exists {
		return nil, errors.New("session not found")
	}

	participant, exists := session.Participants[participantID]
	if !exists {
		return nil, errors.New("participant not found")
	}

	if session.Status != domain.SessionStatusActive {
		return nil, errors.New("session is not active")
	}

	if questionIndex >= len(session.Quiz.Questions) {
		return nil, errors.New("invalid question index")
	}

	// Store answer
	participant.Answers[questionIndex] = answerChoice

	// Check if answer is correct
	question := session.Quiz.Questions[questionIndex]
	correct := false
	for _, choice := range question.Choices {
		if choice.Id == answerChoice && choice.Correct {
			correct = true
			participant.Score += 10 // 10 points per correct answer
			break
		}
	}

	return &domain.AnswerQuestionResponse{
		Success:       true,
		Correct:       correct,
		Score:         participant.Score,
		QuestionIndex: questionIndex,
	}, nil
}

func (s *SessionService) NextQuestion(ctx context.Context, code string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, exists := s.sessions[code]
	if !exists {
		return errors.New("session not found")
	}

	if session.Status != domain.SessionStatusActive {
		return errors.New("session is not active")
	}

	session.CurrentQuestion++
	if session.CurrentQuestion >= len(session.Quiz.Questions) {
		session.Status = domain.SessionStatusFinished
		session.EndTime = time.Now().Unix()
	}

	return nil
}

func (s *SessionService) GetLeaderboard(ctx context.Context, code string) ([]domain.LeaderboardEntry, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	session, exists := s.sessions[code]
	if !exists {
		return nil, errors.New("session not found")
	}

	var participants []*domain.Participant
	for _, p := range session.Participants {
		participants = append(participants, p)
	}

	// Sort by score (descending)
	sort.Slice(participants, func(i, j int) bool {
		return participants[i].Score > participants[j].Score
	})

	var leaderboard []domain.LeaderboardEntry
	for i, p := range participants {
		leaderboard = append(leaderboard, domain.LeaderboardEntry{
			ParticipantID: p.ID,
			Name:          p.Name,
			Score:         p.Score,
			Rank:          i + 1,
		})
	}

	return leaderboard, nil
}

func (s *SessionService) BroadcastToSession(ctx context.Context, code string, packet interface{}) error {
	s.mutex.RLock()
	session, exists := s.sessions[code]
	s.mutex.RUnlock()

	if !exists {
		return errors.New("session not found")
	}

	// Broadcast to all participants
	for _, participant := range session.Participants {
		if participant.Conn != nil {
			// Send packet (implementation depends on NetService)
			// This will be handled by the WebSocket handler
		}
	}

	return nil
}
