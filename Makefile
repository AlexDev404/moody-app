# Makefile

# Variables
SHELL := /bin/bash
.DEFAULT_GOAL := all

# Directories
SRC_DIR := src
SRC_WASM_DIR := src-wasm
BUILD_DIR := build
BIN_DIR := bin
WASM_DIR := $(SRC_DIR)/static/wasm/bundle

# Compiler and flags
GO := go
GOARGS := build
GOFLAGS := 

# Targets
all: build
.PHONY: all initialize build-wasm build-web clean

initialize:
	cd $(SRC_DIR) && $(GO) mod download
	cd $(SRC_WASM_DIR) && $(GO) mod download
	
	mkdir -p $(BIN_DIR)
	mkdir -p $(BUILD_DIR)
	mkdir -p $(WASM_DIR)

build-wasm: initialize
	cd $(SRC_WASM_DIR) && GOOS=js GOARCH=wasm $(GO) $(GOARGS) -o $(BUILD_DIR)/main.wasm

copy-wasm: build-wasm
# cp $(SRC_WASM_DIR)/$(BUILD_DIR)/main.wasm $(BIN_DIR)/main.wasm
# cp "$(SRC_WASM_DIR)/wasm_exec.js" $(BIN_DIR)/wasm_exec.js
	cp $(SRC_WASM_DIR)/$(BUILD_DIR)/main.wasm $(WASM_DIR)/main.wasm
	cp "$(SRC_WASM_DIR)/wasm_exec.js" $(WASM_DIR)/wasm_exec.js

build-web: copy-wasm
	cd $(SRC_DIR) && $(GO) $(GOARGS) -o ../$(BIN_DIR)/main

run: build-web
	$(SHELL) -c "cd $(SRC_DIR) && npm run gow -- main.go"

clean:
	rm -rf $(BUILD_DIR)
	rm -rf $(BIN_DIR)
	rm -rf $(WASM_DIR)
	rm -rf $(SRC_WASM_DIR)/$(BUILD_DIR)
