package service

import (
	"context"
	"encoding/json"
	"fmt"
	"quiz-svc/internal/domain"
	"quiz-svc/internal/infra/message_broker"
	"time"
)

type EventService struct {
	publisher *message_broker.Publisher
}

func NewEventService(publisher *message_broker.Publisher) *EventService {
	return &EventService{
		publisher: publisher,
	}
}

// PublishLeaderboardUpdated publishes leaderboard update event
func (e *EventService) PublishLeaderboardUpdated(ctx context.Context, sessionCode, participantID, participantName string, score, rank int, leaderboard []domain.LeaderboardEntry) error {
	event := domain.LeaderboardUpdatedEvent{
		BaseEvent: domain.BaseEvent{
			Type:      domain.EventTypeLeaderboardUpdated,
			SessionID: sessionCode,
			Timestamp: time.Now(),
			Source:    "quiz-service",
		},
		Data: domain.LeaderboardUpdatedData{
			ParticipantID:   participantID,
			ParticipantName: participantName,
			Score:           score,
			Rank:            rank,
			Leaderboard:     leaderboard,
		},
	}

	return e.publishEvent(ctx, "leaderboard.updated", event)
}

// PublishParticipantJoined publishes participant joined event
func (e *EventService) PublishParticipantJoined(ctx context.Context, sessionCode, participantID, participantName string, totalParticipants int) error {
	event := domain.ParticipantJoinedEvent{
		BaseEvent: domain.BaseEvent{
			Type:      domain.EventTypeParticipantJoined,
			SessionID: sessionCode,
			Timestamp: time.Now(),
			Source:    "quiz-service",
		},
		Data: domain.ParticipantJoinedData{
			ParticipantID:     participantID,
			ParticipantName:   participantName,
			TotalParticipants: totalParticipants,
		},
	}

	return e.publishEvent(ctx, "participant.joined", event)
}

// PublishParticipantLeft publishes participant left event
func (e *EventService) PublishParticipantLeft(ctx context.Context, sessionCode, participantID, participantName string, totalParticipants int) error {
	event := domain.ParticipantLeftEvent{
		BaseEvent: domain.BaseEvent{
			Type:      domain.EventTypeParticipantLeft,
			SessionID: sessionCode,
			Timestamp: time.Now(),
			Source:    "quiz-service",
		},
		Data: domain.ParticipantLeftData{
			ParticipantID:     participantID,
			ParticipantName:   participantName,
			TotalParticipants: totalParticipants,
		},
	}

	return e.publishEvent(ctx, "participant.left", event)
}

// PublishSessionCreated publishes session created event
func (e *EventService) PublishSessionCreated(ctx context.Context, sessionCode, quizID string, quiz *domain.Quiz) error {
	event := domain.SessionCreatedEvent{
		BaseEvent: domain.BaseEvent{
			Type:      domain.EventTypeSessionCreated,
			SessionID: sessionCode,
			Timestamp: time.Now(),
			Source:    "quiz-service",
		},
		Data: domain.SessionCreatedData{
			Code:   sessionCode,
			QuizID: quizID,
			Quiz:   quiz,
		},
	}

	return e.publishEvent(ctx, "session.created", event)
}

// PublishSessionStarted publishes session started event
func (e *EventService) PublishSessionStarted(ctx context.Context, sessionCode string, totalParticipants int, quiz *domain.Quiz) error {
	event := domain.SessionStartedEvent{
		BaseEvent: domain.BaseEvent{
			Type:      domain.EventTypeSessionStarted,
			SessionID: sessionCode,
			Timestamp: time.Now(),
			Source:    "quiz-service",
		},
		Data: domain.SessionStartedData{
			TotalParticipants: totalParticipants,
			Quiz:              quiz,
		},
	}

	return e.publishEvent(ctx, "session.started", event)
}

// PublishSessionEnded publishes session ended event
func (e *EventService) PublishSessionEnded(ctx context.Context, sessionCode string, finalLeaderboard []domain.LeaderboardEntry, totalParticipants int) error {
	event := domain.SessionEndedEvent{
		BaseEvent: domain.BaseEvent{
			Type:      domain.EventTypeSessionEnded,
			SessionID: sessionCode,
			Timestamp: time.Now(),
			Source:    "quiz-service",
		},
		Data: domain.SessionEndedData{
			FinalLeaderboard:  finalLeaderboard,
			TotalParticipants: totalParticipants,
		},
	}

	return e.publishEvent(ctx, "session.ended", event)
}

// PublishAnswerSubmitted publishes answer submitted event
func (e *EventService) PublishAnswerSubmitted(ctx context.Context, sessionCode, participantID, participantName string, questionIndex int, answerChoice string, correct bool, score, newTotalScore int) error {
	event := domain.AnswerSubmittedEvent{
		BaseEvent: domain.BaseEvent{
			Type:      domain.EventTypeAnswerSubmitted,
			SessionID: sessionCode,
			Timestamp: time.Now(),
			Source:    "quiz-service",
		},
		Data: domain.AnswerSubmittedData{
			ParticipantID:   participantID,
			ParticipantName: participantName,
			QuestionIndex:   questionIndex,
			AnswerChoice:    answerChoice,
			Correct:         correct,
			Score:           score,
			NewTotalScore:   newTotalScore,
		},
	}

	return e.publishEvent(ctx, "answer.submitted", event)
}

// publishEvent is a helper method to publish events
func (e *EventService) publishEvent(ctx context.Context, subject string, event interface{}) error {
	if e.publisher == nil {
		return fmt.Errorf("publisher not initialized")
	}

	// Convert event to JSON
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publish event
	return e.publisher.Publish(subject, event.(domain.BaseEvent).Type, eventData)
}
