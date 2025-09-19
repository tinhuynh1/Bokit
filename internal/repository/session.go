package repository

import (
	"context"
	"quiz-svc/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type sessionCollection struct {
	collection *mongo.Collection
}

func NewSessionCollection(db *mongo.Database) domain.SessionCollection {
	collection := db.Collection("sessions")
	return &sessionCollection{collection: collection}
}

func (r *sessionCollection) CreateSession(ctx context.Context, session *domain.SessionCreate) error {
	_, err := r.collection.InsertOne(ctx, session)
	return err
}

func (r *sessionCollection) GetSessionByCode(ctx context.Context, code string) (*domain.Session, error) {
	var session domain.Session
	err := r.collection.FindOne(ctx, bson.M{"code": code}).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &session, err
}
