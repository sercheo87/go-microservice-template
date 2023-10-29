# Go compiler
GO := go

# Output binary name
BINARY_NAME := go-microservice-template

# Directories
BIN_DIR := out
DEPLOYMENT_DIR := deployments
DOCS_DOCKER_IMAGE = asciidoctor/docker-asciidoctor
DOCS_SOURCE_DIR = docs/asciidoc
DOCS_IMAGES_DIR = docs/asciidoc/images
DOCS_OUTPUT_DIR = out/docs

all: docker_build

docker_build:
	@echo "Building..."
	@docker build -t go-microservice-template:1.0.0 -f Dockerfile .

.PHONY: docker_build
