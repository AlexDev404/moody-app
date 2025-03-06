# Makefile

# Variables
SHELL := /bin/bash
.DEFAULT_GOAL := all

# Directories
SRC_DIR := src
BUILD_DIR := build
BIN_DIR := bin

# Compiler and flags
GO := go build
GOFLAGS := 

# Targets
all: build-wasm
.PHONY: all build clean

initialize:
	mkdir -p $(BIN_DIR)
	mkdir -p $(BUILD_DIR)

build-wasm:
	GOOS=js GOARCH=wasm $(GO) -o $(BUILD_DIR)/main.wasm

build: $(BUILD_DIR)/main.o
	$(GO) $(GOFLAGS) -o $(BIN_DIR)/main $(BUILD_DIR)/main.o

$(BUILD_DIR)/main.wasm: build
	mkdir -p $(BUILD_DIR)
	GOOS=js GOARCH=wasm $(GO) build -o $(BUILD_DIR)/main.wasm $(SRC_DIR)/main.go

$(BUILD_DIR)/main.o: $(SRC_DIR)/main.c
	mkdir -p $(BUILD_DIR)
	$(GO) $(GOFLAGS) -c $(SRC_DIR)/main.c -o $(BUILD_DIR)/main.o

clean:
	rm -rf $(BUILD_DIR) $(BIN_DIR)
