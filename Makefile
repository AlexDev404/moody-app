# Makefile
include .envrc

# Variables
SHELL := /bin/bash
WSL_CHECK :=
.DEFAULT_GOAL := all

# Directories
SRC_DIR := src
SRC_WASM_DIR := src-wasm
BUILD_DIR := build
BIN_DIR := bin
WASM_DIR := $(SRC_DIR)/static/wasm/bundle
WASM_DIR_WIN := $(SRC_DIR)\static\wasm\bundle

# Compiler and flags
GO := go
GOARGS := build
GOFLAGS := 

# Targets
all: run
.PHONY: all initialize build-web build-wasm copy-wasm clean create_migrations prepare

# ------------------ BEGIN PLATFORM AND ARCHITECTURE DETECTION --------------------
BUILD_PLATFORM=
BUILD_ARCH=

ifeq ($(OS),Windows_NT)
	BUILD_PLATFORM = WIN32
	ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
		BUILD_ARCH = AMD64
	else
	ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
		BUILD_ARCH = AMD64
	endif
	ifeq ($(PROCESSOR_ARCHITECTURE),x86)
		BUILD_ARCH = IA32
	endif
endif
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		BUILD_PLATFORM = LINUX
		WSL_CHECK = $(shell if [ -f /proc/sys/fs/binfmt_misc/WSLInterop ]; then echo true; else echo false; fi)
	endif
	ifeq ($(UNAME_S),Darwin)
		BUILD_PLATFORM = OSX
	endif
	UNAME_P := $(shell uname -p)
	ifeq ($(UNAME_P),x86_64)
		BUILD_ARCH = AMD64
	endif
	ifneq ($(filter %86,$(UNAME_P)),)
		BUILD_ARCH = IA32
	endif
	ifneq ($(filter arm%,$(UNAME_P)),)
		BUILD_ARCH = ARM
	endif
	UNAME_P = $(shell uname -m)
	ifeq ($(UNAME_P),x86_64)
		BUILD_ARCH = AMD64
	endif
	ifneq ($(filter %86,$(UNAME_P)),)
		BUILD_ARCH = IA32
	endif
endif
# --------------- END PLATFORM AND ARCHITECTURE DETECTION -----------------
prepare:
	cd $(SRC_DIR) && npm install
ifeq ($(BUILD_PLATFORM),LINUX)
	$(GO) install github.com/mitranim/gow@latest
endif

initialize:
	@echo "The operating system of this computer is: $(BUILD_PLATFORM)"
	@echo "The processor architecture of this system is: $(BUILD_ARCH)"
	cd $(SRC_DIR) && $(GO) mod download
	cd $(SRC_WASM_DIR) && $(GO) mod download
	
ifeq ($(BUILD_PLATFORM),LINUX)
		mkdir -p $(BIN_DIR)
		mkdir -p $(BUILD_DIR)
		mkdir -p $(WASM_DIR)
else
ifeq ($(BUILD_PLATFORM),WIN32)
		powershell.exe -Command "foreach ($$dir in @('$(BIN_DIR)', '$(BUILD_DIR)')) { if (!(Test-Path -Path $$dir)) { New-Item -ItemType Directory -Path $$dir -Force } }"
endif
endif

build-web: initialize
	cd $(SRC_DIR) && $(GO) $(GOARGS) -o ../$(BIN_DIR)/main

# Define this variable at the top level
run: build-web copy-wasm
ifeq ($(BUILD_PLATFORM),LINUX)
# WSL detection
ifeq ($(WSL_CHECK),true)
	cd $(SRC_DIR) && npm run dev -- . --dsn ${DB_DSN} --openai-key ${OPENAI_API_KEY} --jwt-secret ${JWT_SECRET}
else
	cd $(SRC_DIR) && npm run gow -- . --dsn ${DB_DSN} --openai-key ${OPENAI_API_KEY} --jwt-secret ${JWT_SECRET}
endif

else
ifeq ($(BUILD_PLATFORM),WIN32)
	cd $(SRC_DIR) && go run . --openai-key ${OPENAI_API_KEY} --dsn ${DB_DSN} --jwt-secret ${JWT_SECRET}
endif
endif

build-wasm: build-web
ifeq ($(BUILD_PLATFORM),LINUX)
	cd $(SRC_WASM_DIR) && GOOS=js GOARCH=wasm $(GO) $(GOARGS) -o $(BUILD_DIR)/main.wasm
else
ifeq ($(BUILD_PLATFORM),WIN32)
	cd $(SRC_WASM_DIR) && powershell.exe -Command "$$env:GOOS='js'; $$env:GOARCH='wasm'; & $(GO) $(GOARGS) -o $(BUILD_DIR)/main.wasm"
endif
endif

copy-wasm: build-wasm
# cp $(SRC_WASM_DIR)/$(BUILD_DIR)/main.wasm $(BIN_DIR)/main.wasm
# cp "$(SRC_WASM_DIR)/wasm_exec.js" $(BIN_DIR)/wasm_exec.js
ifeq ($(BUILD_PLATFORM),LINUX)
	cp $(SRC_WASM_DIR)/$(BUILD_DIR)/main.wasm $(WASM_DIR)/main.wasm
	cp "$(SRC_WASM_DIR)/wasm_exec.js" $(WASM_DIR)/wasm_exec.js
else
ifeq ($(BUILD_PLATFORM),WIN32)
	copy $(SRC_WASM_DIR)\$(BUILD_DIR)\main.wasm $(WASM_DIR_WIN)\main.wasm
	copy "$(SRC_WASM_DIR)\wasm_exec.js" $(WASM_DIR_WIN)\wasm_exec.js
endif
endif

clean:
ifeq ($(BUILD_PLATFORM),LINUX)
	rm -rf $(BUILD_DIR)
	rm -rf $(BIN_DIR)
else
ifeq ($(BUILD_PLATFORM),WIN32)
	rd /s /q $(BUILD_DIR)
	rd /s /q $(BIN_DIR)
endif
endif

create_migrations:
	migrate create -seq -ext=.sql -dir=./migrations create_feedback_table

migrate:
ifeq ($(BUILD_PLATFORM),LINUX)
	migrate -path ./migrations -database ${DB_DSN} up
else
ifeq ($(BUILD_PLATFORM),WIN32)
	migrate -path ./migrations -database ${DB_DSN} up
endif
endif

db/psql:
	psql ${DB_DSN}

dumb_migrations:
	migrate create -seq -ext=.sql -dir=./migrations create {something} table