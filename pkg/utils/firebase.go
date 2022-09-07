package utils

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FCMClient struct {
	Client *messaging.Client
}

func NewFirebaseService(ctx context.Context) (*firebase.App, error) {
	// decodedKey, err := getDecodedFireBaseKey()
	// if err != nil {
	// 	return nil, err
	// }

	opts := []option.ClientOption{option.WithCredentialsFile("./mana-notification-service.json")}

	firebaseApp, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		return nil, err
	}

	return firebaseApp, nil
}

func NewFCMClient(ctx context.Context) (*FCMClient, error) {
	firebaseApp, err := NewFirebaseService(ctx)
	if err != nil {
		return nil, err
	}

	client, err := firebaseApp.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	return &FCMClient{
		Client: client,
	}, nil
}

func (c *FCMClient) PushNotification(title, body, firebaseToken string) (string, error) {
	return c.Client.Send(context.Background(), &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Token: firebaseToken,
	})
}
