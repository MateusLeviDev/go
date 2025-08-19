# Variables
PROJECT_ID := demo-test
PUBSUB_EMULATOR_HOST := localhost:8085
FUNCTIONS_DIR := functions
WEB_DIR := cmd/main.go

# Start the Pub/Sub emulator
# emulator:
# 	gcloud beta emulators pubsub start --project=$(PROJECT_ID) --host-port=$(PUBSUB_EMULATOR_HOST) --quiet &

# # Create Pub/Sub topic and subscription
# setup-pubsub:
# 	export PUBSUB_EMULATOR_HOST=$(PUBSUB_EMULATOR_HOST); \
# 	gcloud pubsub topics create send-event-topic --project=$(PROJECT_ID); \
# 	gcloud pubsub subscriptions create send-event-sub --project=$(PROJECT_ID) --topic=send-event-topic --push-endpoint=http://localhost:8081/

emulator:
	docker compose up -d pubsub-emulator

stop-emulator:
	docker compose down

logs-emulator:
	docker compose logs -f pubsub-emulator

# Run the Cloud Function locally
run-function:
	cd $(FUNCTIONS_DIR) && \
	export PUBSUB_EMULATOR_HOST=$(PUBSUB_EMULATOR_HOST); \
	export PUBSUB_PROJECT_ID=$(PROJECT_ID); \
	go run ./main.go

# Run the web server
run-web:
	export PUBSUB_EMULATOR_HOST=$(PUBSUB_EMULATOR_HOST); \
	export PUBSUB_PROJECT_ID=$(PROJECT_ID); \
	go run ./$(WEB_DIR) -gcp-project-id=$(PROJECT_ID)

# Run all components in the background
run:
	make emulator
	@sleep 2 # Wait for emulator to start
	make setup-pubsub
	make run-function &
	make run-web &

# Clean up
clean:
	-pkill -f "pubsub" # Stop the Pub/Sub emulator
	rm -rf firebase-data

.PHONY: emulator setup-pubsub run-function run-web run clean