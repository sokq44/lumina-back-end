SRC_DIR=/home/sokq/Projects/lumina/lumina-back-end/src
BIN_DIR=/home/sokq/Projects/lumina/lumina-back-end/bin
APP=$(BIN_DIR)/lumina-backend
USER=deployer
HOST=161.35.199.143
DIR=/home/deployer/lumina-backend

all: build

dev:
	@bash -c 'export $$(grep -v "^#" .env | xargs); clear; cd src && go run . -v=true -d=true -l ../logs'


build:
	@mkdir -p $(BIN_DIR)
	@printf "Building the application..."
	@cd $(SRC_DIR) || exit; go build -o $(APP) main.go
	@printf "\r"
	@printf "\rApplication has been built.\n"

deploy: build
	@printf "Stopping the application service..."
	@ssh $(USER)@$(HOST) 'sudo systemctl stop lumina-backend.service'
	@printf "\r"
	@printf "Application service has been stopped.\nSending files..."
	@scp -q $(APP) $(USER)@$(HOST):$(DIR)
	@printf "\r"
	@printf "Files have been sent.\nRestarting the application service..."
	@ssh $(USER)@$(HOST) 'sudo systemctl start lumina-backend.service'
	@printf "\r"
	@printf "Application service has been restarted.\n"