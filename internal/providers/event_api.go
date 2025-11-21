package providers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/abkristanto/go-microservice/internal/models"
	"github.com/abkristanto/go-microservice/internal/providers/dtos"
	"github.com/abkristanto/go-microservice/internal/providers/mappers"
)

type EventProvider interface {
	GetEvents() ([]models.Event, error)
	APISource() string
}

type HTTPEventProvider struct {
	baseURL string
	client  *http.Client
	source  string
}

func NewHTTPEventProvider(baseURL string) EventProvider {
	return &HTTPEventProvider{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		source: baseURL,
	}
}

func (p *HTTPEventProvider) GetEvents() ([]models.Event, error) {
	req, err := http.NewRequest(http.MethodGet, p.baseURL+"/events", nil)
	if err != nil {
		log.Printf("GetEvents: building request failed: %v", err)
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		log.Printf("GetEvents: request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("GetEvents: unexpected status code: %d", resp.StatusCode)
		return nil, err
	}

	var events []dtos.Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		log.Printf("GetEvents: decode failed: %v", err)
		return nil, err
	}

	eventsModels := make([]models.Event, 0, len(events))

	for _, e := range events {
		em := mappers.ToDomainEvent(e)
		eventsModels = append(eventsModels, em)
	}

	return eventsModels, nil
}

func (p *HTTPEventProvider) APISource() string {
	return p.source
}
