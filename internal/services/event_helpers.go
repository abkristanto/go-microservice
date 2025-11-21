package services

import "github.com/abkristanto/go-microservice/internal/models"

type payload struct {
	ChangeType       string
	APISource        string
	ResourceLocation string
	Event            models.Event
}

func hasChanged(stored models.Event, remote models.Event) bool {

	if stored.Title != remote.Title {
		return true
	}
	if stored.Description != remote.Description {
		return true
	}
	if !stored.StartsAt.Equal(remote.StartsAt) {
		return true
	}

	return false
}