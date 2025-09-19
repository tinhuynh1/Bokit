package service

import (
	"context"
	"quiz-svc/internal/domain"

	"quiz-svc/internal/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type QuizService struct {
	logger         *zap.Logger
	quizCollection domain.QuizCollection
}

func NewQuizService(logger *zap.Logger, quizCollection domain.QuizCollection) *QuizService {
	return &QuizService{logger: logger, quizCollection: quizCollection}
}

func (s *QuizService) ListQuiz(ctx context.Context) ([]domain.Quiz, error) {
	quizzes, err := s.quizCollection.GetQuizzes(ctx)
	if err != nil {
		s.logger.Error("failed to get quizzes", zap.Error(err))
		return nil, err
	}
	return quizzes, nil
}

func (s *QuizService) GetQuizById(ctx context.Context, id string) (*domain.Quiz, error) {
	quiz, err := s.quizCollection.GetQuizById(ctx, id)
	if err != nil {
		s.logger.Error("failed to get quiz by id", zap.Error(err))
		return nil, err
	}
	return quiz, nil
}

func (s *QuizService) UpdateQuiz(ctx context.Context, quiz *domain.Quiz) error {
	err := s.quizCollection.UpdateQuiz(ctx, quiz)
	if err != nil {
		s.logger.Error("failed to update quiz", zap.Error(err))
		return err
	}
	return nil
}

func (s *QuizService) CreateQuiz(ctx context.Context, quiz *dto.CreateQuizRequest) error {
	s.logger.Info("Creating quiz", zap.String("name", quiz.Name))
	quizDomain := &domain.Quiz{
		Name:      quiz.Name,
		Questions: make([]domain.QuizQuestion, len(quiz.Questions)),
	}
	for i, question := range quiz.Questions {
		quizDomain.Questions[i] = domain.QuizQuestion{
			Id:      primitive.NewObjectID(),
			Name:    question.Name,
			Time:    question.Time,
			Choices: make([]domain.QuizChoice, len(question.Choices)),
		}
		for j, choice := range question.Choices {
			quizDomain.Questions[i].Choices[j] = domain.QuizChoice{
				Id:      primitive.NewObjectID(),
				Name:    choice.Name,
				Correct: choice.Correct,
			}
		}
	}
	err := s.quizCollection.CreateQuiz(ctx, quizDomain)
	if err != nil {
		s.logger.Error("❌ Failed to create quiz",
			zap.String("quiz_name", quiz.Name),
			zap.String("error", err.Error()))
		return err
	}

	s.logger.Info("Quiz created successfully",
		zap.String("quiz_name", quiz.Name),
		zap.String("quiz_id", quizDomain.Id.Hex()))
	return nil
}
