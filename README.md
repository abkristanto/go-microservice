# Go Event Sync Microservice

This project is an end-to-end microservice that:
- Polls an external REST API for Event data
- Upserts events into MongoDB
- Detects created / updated / deleted events
- Writes what changed to an outbox table
- Runs a worker that reads the outbox and publishes messages to a broker (currently implemented as a LoggingProducer that just prints the payload)

The design of the microservice uses the transactional outbox pattern:

> HTTP -> Service (sync + diff) -> Repository (Mongo + transactions) -> Outbox -> Worker -> Broker

## High Level Flow
1. Sync worker wakes up on an interval of 30s, fetching events from the external API via the event provider
2. For each event it checks:
    - If it only exists remotely -> event.created
    - If it exists in both places but fields differ -> event.updated
    - If it exists only in DB -> event.deleted
3. For each change, the service:
    - Upserts the event in the database within a transaction
    - Writes an outbox row describing the change
4. Outbox worker periodically reads the outbox rows with pending status, and publishes the payload, finally marking them as sent.

## Tech Stack
- Language: Go
- Database: MongoDB (single node, replica set enabled)
- External API: FastAPI exposing `/events`
- Containerization: Docker
- Messaging: Outbox + Producer interface

## How To Run

To start the whole stack (MongoDB, the Go sync-service, and the FastAPI events-api), just run:

```bash
make dev
```

This command will:
1. Start MongoDB in Docker
2. Wait a few seconds, then initialize a single-node replica set so MongoDB transactions work
3. Build and start `events-api`, the FastAPI app that serves `GET /events` and `sync-service` which is the Go microservice that polls `events-api`, upserts into Mongo, writes change events into the outbox, and publishes them via a logging producer. The FastAPI endpoint returns events with a start time that differ each call; so you can directly visualize how the service works when handling differing fields.
4. Attach to the logs for `sync-service` and `events-api` so you can watch activity in real time.

### Example logs

```logs
events-api    | INFO:     172.23.0.4:47478 - "GET /events HTTP/1.1" 200 OK
sync-service  | 2025/11/21 03:57:41 SyncEvents completed successfully took=31.175625ms
sync-service  | 2025/11/21 03:57:46 Publishing message to broker: {"ChangeType":"event.updated","APISource":"http://events-api:8000","ResourceLocation":"events_db.events.0a170740-d235-4514-bafb-7124cf7359b6","Event":{"ID":"0a170740-d235-4514-bafb-7124cf7359b6","ExternalID":"evt_1","Title":"Test Event 1","Description":"First test event","StartsAt":"2025-11-21T03:57:41.454777Z"}}
sync-service  | 2025/11/21 03:57:46 Publishing message to broker: {"ChangeType":"event.updated","APISource":"http://events-api:8000","ResourceLocation":"events_db.events.9f0eb01e-9704-458f-b683-e2ad2c23ae2d","Event":{"ID":"9f0eb01e-9704-458f-b683-e2ad2c23ae2d","ExternalID":"evt_2","Title":"Test Event 2","Description":"Second test event","StartsAt":"2025-11-21T04:57:41.454777Z"}}
```

### Stopping the stack
When you're done, stop and remove the containers with:
```bash
make down
```
This ensures that all Docker containers started by the project are fully stopped and cleaned up,

## Project Structure

```text
├── Dockerfile                      # Go sync-service Dockerfile
├── Makefile                        # dev, down, logs targets
├── README.md
├── cmd
│   └── service
│       └── main.go                 # Entry point: wiring, DI, workers
├── docker-compose.yml              # Orchestrates mongo, events-api, sync-service
├── fastapi_app
│   ├── Dockerfile                  # FastAPI Dockerfile
│   └── main.py                     # /events provider for local dev
├── go.mod
├── go.sum
└── internal
    ├── mocks                       # Simple in-memory mocks for interfaces
    │   └── mock.go
    ├── models                      # Domain models
    │   ├── event.go                # Event domain model
    │   └── outbox.go               # Outbox domain model
    ├── providers                   # External API integration
    │   ├── dtos
    │   │   └── event.go            # JSON shape from /events
    │   ├── event_api.go            # HTTPEventProvider implementation
    │   └── mappers
    │       └── event_mapper.go     # DTO -> domain mapper
    ├── repositories                # Persistence interfaces + Mongo adapters
    │   ├── event_repository.go     # EventRepository interface
    │   ├── mongo
    │   │   ├── documents           # Mongo document types
    │   │   │   ├── event_document.go
    │   │   │   └── outbox_document.go
    │   │   ├── mappers             # domain <-> document mappers
    │   │   │   ├── event_mapper.go
    │   │   │   └── outbox_mapper.go
    │   │   ├── mongo_event_repository.go
    │   │   ├── mongo_outbox_repository.go
    │   │   └── mongo_transaction_manager.go
    │   ├── outbox_repository.go    # OutboxRepository interface
    │   └── transaction_manager.go  # TransactionManager interface
    ├── services                    # Application / domain logic
    │   ├── event_helpers.go        # Payload struct, hasChanged, helpers
    │   ├── event_service.go        # SyncEvents implementation
    │   └── logging.go              # Logging decorator for EventService
    └── workers                     # Background workers
        ├── outbox_worker.go        # Reads outbox + calls Producer
        └── sync_worker.go          # Triggers SyncEvents on an interval

```