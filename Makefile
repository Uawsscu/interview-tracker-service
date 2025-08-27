###################
# Docker commands
up:
	docker-compose up -d

down:
	docker-compose down --remove-orphans

clean:
	sudo rm -rf db/data

###################
# Run Main and Swagger
.PHONY: swagger-clean swagger-up

swagger-clean:
	@rm -f docs/docs.go docs/swagger.json docs/swagger.yaml

swagger-up: swagger-clean
	@export PATH=$$(go env GOPATH)/bin:$$PATH; \
	echo ">> generating swagger..."; \
	swag init -d ./cmd/api,./internal -g main.go -o ./docs --parseInternal

###################
# unit tests
########################
# Go Test + Coverage
########################
SHELL := /bin/bash

# โฟลเดอร์เก็บผล coverage
COVER_DIR := coverage
COVER_PROFILE := $(COVER_DIR)/coverage.out
COVER_HTML := $(COVER_DIR)/coverage.html

# เลือกแพ็กเกจทั้งหมด ยกเว้นโฟลเดอร์ mock (ถ้าต้องการครอบคลุมทั้งหมด ให้ใช้ go list ./... ตรง ๆ)
PKGS := $(shell go list ./... | grep -v /internal/mocks | grep -v /mock_repo)

# ใช้สำหรับ -coverpkg ให้คิด coverage แบบข้ามแพ็กเกจ
COVERPKG := $(shell go list ./... | tr '\n' ',' | sed 's/,$$//')

.PHONY: test testv cover cover-html cover-func cover-each clean \
        mockgen-install gen

## รันเทสต์ทั้งหมด + สร้าง coverage profile รวม
test:
	@mkdir -p $(COVER_DIR)
	@echo ">> running tests with coverage..."
	@go test -race -covermode=atomic -coverpkg=$(COVERPKG) -coverprofile=$(COVER_PROFILE) $(PKGS)

## รันเทสต์แบบ verbose
testv:
	@mkdir -p $(COVER_DIR)
	@echo ">> running verbose tests..."
	@go test -v -race -covermode=atomic -coverpkg=$(COVERPKG) -coverprofile=$(COVER_PROFILE) $(PKGS)

## แสดงสรุป % coverage ในเทอร์มินัล
cover: test
	@echo ">> coverage summary"
	@go tool cover -func=$(COVER_PROFILE) | tail -n 1

## สร้าง HTML รายงาน coverage
cover-html: test
	@echo ">> generating HTML report at $(COVER_HTML)"
	@go tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)
	@echo "Open file://$(PWD)/$(COVER_HTML)"

## แยก coverage ต่อแพ็กเกจ (มีไฟล์ย่อยใน coverage/pkg/)
cover-each:
	@rm -rf $(COVER_DIR)/pkg && mkdir -p $(COVER_DIR)/pkg
	@set -e; \
	for pkg in $(PKGS); do \
	  out="$(COVER_DIR)/pkg/$$(echo $$pkg | tr '/' '_').out"; \
	  echo ">> $$pkg"; \
	  go test -race -covermode=atomic -coverpkg=$(COVERPKG) -coverprofile="$$out" "$$pkg"; \
	done

## ล้างไฟล์ผลลัพธ์
clean:
	@rm -rf $(COVER_DIR)

################### SETUP MOCKGEN ##############################
setup-mockgen:
	@echo "Setting up mockgen"
	@go install github.com/golang/mock/mockgen@v1.6.0
	@export PATH=$$PATH:$(shell go env GOPATH)/bin

################### GENERATED MOCK REPO ########################
INTERFACE_DIR_REPO := ./internal/adapters/repositories
MOCK_DIR_REPO := ./internal/mocks/mock_repo
MOCKGEN := mockgen
INTERFACE_FILES_REPO := $(wildcard $(INTERFACE_DIR_REPO)/*.go)
MOCK_FILES_REPO := $(patsubst $(INTERFACE_DIR_REPO)/%.go,$(MOCK_DIR_REPO)/mock_%.go,$(INTERFACE_FILES_REPO))
mock-repo: $(MOCK_FILES_REPO)
$(MOCK_DIR_REPO)/mock_%.go: $(INTERFACE_DIR_REPO)/%.go
	@mkdir -p $(MOCK_DIR_REPO)
	$(MOCKGEN) -source=$< -destination=$@ -package=mock_repo
clean-repo:
	rm -rf $(MOCK_DIR_REPO)/*.go
.PHONY: mock-repo clean-repo