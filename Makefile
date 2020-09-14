.PHONY: build
build:
	go build -v ./cmd/chat

.DEFAULT_GOAL := build