COMPOSE := docker compose

.PHONY: dev down logs

dev:
	$(COMPOSE) up -d mongo
	# wait a bit for mongod to come up
	sleep 5
	# init replica set (ignore "already initialized" error)
	$(COMPOSE) exec -T mongo mongosh --eval 'rs.initiate({_id: "rs0", members: [{ _id: 0, host: "mongo:27017" }]})' || true
	$(COMPOSE) up -d --build events-api sync-service
	$(COMPOSE) logs -f sync-service events-api

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f
