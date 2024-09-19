#!/bin/env bash

echo "Building webhook listener..."
go build -o webhook_listener ./cmd/webhook_listener

echo "Moving binary to /usr/local/bin..."
sudo mv webhook_listener /usr/local/bin/

echo "Setting up systemd service..."
sudo cp service/webhook_listener.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable webhook_listener.service
sudo systemctl start webhook_listener.service

echo "Deployment complete."

