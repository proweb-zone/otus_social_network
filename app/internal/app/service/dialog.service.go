package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"otus_social_network/app/internal/app/repository"

	eventclient "github.com/proweb-zone/event-client"
	pb "github.com/proweb-zone/event-client/gen/go"
)

type DialogService struct {
	repo   *repository.FriendsRepository
	client *eventclient.EventClient
}

func NewDialogService(
	newRepo *repository.FriendsRepository,
	newClient *eventclient.EventClient,
) *DialogService {
	return &DialogService{
		repo:   newRepo,
		client: newClient,
	}
}

type EventReponse struct {
	User_id_sender    int
	User_id_recipient int
	Id                int
	Msg               string
}

func (d *DialogService) CheckUserAccess(event *pb.Event) (*pb.PublishResponse, error) {
	payload := event.GetPayload()

	parsedResponse := &EventReponse{}
	if err := json.Unmarshal(payload, &parsedResponse); err != nil {
		log.Fatalf("Failed to parse payload JSON: %v", err)
	}

	state := true

	userId := parsedResponse.User_id_sender
	friendId := parsedResponse.User_id_recipient
	_, err := d.repo.GetFriendById(userId, friendId)
	if err != nil {
		state = false
	}

	return d.client.Publish(context.Background(), &pb.Event{
		Type:    "user.access",
		Source:  "otus-service",
		Payload: []byte(fmt.Sprintf(`{"id": %d, "state": %t}`, parsedResponse.Id, state)),
	})
}
