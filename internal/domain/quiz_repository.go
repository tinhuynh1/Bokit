package domain

import "context"

type QuizCollection interface {
	CreateQuiz(ctx context.Context, quiz *Quiz) error
	GetQuizById(ctx context.Context, id string) (*Quiz, error)
	GetQuizzes(ctx context.Context) ([]Quiz, error)
	UpdateQuiz(ctx context.Context, quiz *Quiz) error
}
