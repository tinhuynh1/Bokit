package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Quiz struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Questions []QuizQuestion     `json:"questions" bson:"questions"`
}

type QuizQuestion struct {
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name"`
	Time    int                `json:"time" bson:"time"`
	Choices []QuizChoice       `json:"choices" bson:"choices"`
}

type QuizChoice struct {
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name"`
	Correct bool               `json:"correct" bson:"correct"`
}
