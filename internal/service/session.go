package service

import (
	"context"
	"errors"
	"math/rand"
	"quiz-svc/internal/domain"
	"quiz-svc/pkg/logger"
	"strconv"
	"sync"

	"quiz-svc/internal/repository"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type SessionService struct {
	sessions          map[string][]domain.Participant
	sessionCollection domain.SessionCollection
	redisRepo         *repository.RedisRepo
	mutex             sync.RWMutex
}

func NewSessionService(sessionCollection domain.SessionCollection, redisRepo *repository.RedisRepo) *SessionService {
	return &SessionService{
		sessions:          make(map[string][]domain.Participant),
		sessionCollection: sessionCollection,
		redisRepo:         redisRepo,
	}
}

func (s *SessionService) CreateSession(ctx context.Context, quiz *domain.Quiz) error {

	sessionCreate := &domain.SessionCreate{
		Code:   s.generateCode(),
		QuizId: quiz.Id.Hex(),
		Status: domain.SessionStatusWaiting,
	}

	for {
		session, err := s.sessionCollection.GetSessionByCode(ctx, sessionCreate.Code)
		if err != nil {
			return errors.New("failed to get session")
		}
		if session == nil {
			break
		}
		sessionCreate.Code = s.generateCode()
	}

	err := s.sessionCollection.CreateSession(ctx, sessionCreate)
	if err != nil {
		return errors.New("failed to create session")
	}
	//public message session created
	return nil
}

func (s *SessionService) generateCode() string {
	return strconv.Itoa(100000 + rand.Intn(900000))
}

func (s *SessionService) JoinSession(ctx context.Context, code, participantName string, conn *websocket.Conn) (*domain.Participant, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	session, err := s.sessionCollection.GetSessionByCode(ctx, code)
	if err != nil {
		return nil, errors.New("failed to get session")
	}
	if session == nil {
		return nil, errors.New("session not found")
	}

	if session.Status != domain.SessionStatusWaiting {
		return nil, errors.New("session is not accepting new participants")
	}

	participant := &domain.Participant{
		ID:      uuid.New().String(),
		Code:    code,
		Name:    participantName,
		Score:   0,
		Answers: make(map[int]string),
	}
	s.sessions[code] = append(s.sessions[code], *participant)
	err = s.redisRepo.Publish(ctx, "session:participant_joined", participant)
	if err != nil {
		return nil, errors.New("failed to publish participant joined event")
	}
	logger.L.Info("participant joined session", zap.String("session", code), zap.String("participant", participantName))
	return participant, nil
}

func (s *SessionService) HandleParticipantJoined(ctx context.Context) {
	message, err := s.redisRepo.Subscribe(ctx, "session:participant_joined")
	if err != nil {
		logger.L.Error("failed to subscribe to participant joined event", zap.Error(err))
		return
	}
	logger.L.Info("participant joined session", zap.String("message", message))
}
