dev:
	make dev:chat

dev\:chat:
	cd chat && go build -o ./chat && ./chat -host=":8080"