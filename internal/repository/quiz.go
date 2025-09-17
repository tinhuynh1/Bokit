package repository

import (
	"context"
	"fmt"
	"quiz-svc/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type quizCollection struct {
	collection *mongo.Collection
}

func NewQuizCollection(db *mongo.Database) domain.QuizCollection {
	collection := db.Collection("quizzes")
	return &quizCollection{collection: collection}
}

func (r *quizCollection) CreateQuiz(ctx context.Context, quiz *domain.Quiz) error {
	_, err := r.collection.InsertOne(ctx, quiz)
	return err
}

func (r *quizCollection) GetQuizById(ctx context.Context, id string) (*domain.Quiz, error) {
	var quiz domain.Quiz
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error converting id to object id", err)
		return nil, err
	}
	err = r.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&quiz)
	if err != nil {
		fmt.Println("Error getting quiz by id", err)
		return nil, err
	}
	return &quiz, nil
}

func (r *quizCollection) GetQuizzes(ctx context.Context) ([]domain.Quiz, error) {
	var quizzes []domain.Quiz
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &quizzes)
	if err != nil {
		return nil, err
	}
	return quizzes, nil
}

func (r *quizCollection) UpdateQuiz(ctx context.Context, quiz *domain.Quiz) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{
		"_id": quiz.Id,
	}, bson.M{
		"$set": quiz,
	})

	return err
}
