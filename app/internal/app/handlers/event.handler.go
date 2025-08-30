package handlers

import (
	"log"

	pb "github.com/proweb-zone/event-client/gen/go"
)

func (h *Handler) Events(event *pb.Event) error {
	switch event.Type {
	case "dialog.send":
		h.dialogService.CheckUserAccess(event)
		return nil
	default:
		log.Printf("Unknown event type: %s", event.Type)
		return nil
	}
}
