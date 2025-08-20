package google

import (
	"context"
	"errors"

	"google.golang.org/api/idtoken"
)

type Client struct {
	Audience string // Client ID của Google OAuth
}

type GoogleProfile struct {
	Email string
	Name  string
}

func NewClient(audience string) *Client {
	return &Client{
		Audience: audience,
	}
}

func (c *Client) VerifyIDToken(ctx context.Context, token string) (*GoogleProfile, error) {
	payload, err := idtoken.Validate(ctx, token, c.Audience)
	if err != nil {
		return nil, err
	}

	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return nil, errors.New("invalid email in Google token")
	}

	name, _ := payload.Claims["name"].(string) // optional

	return &GoogleProfile{
		Email: email,
		Name:  name,
	}, nil
}
