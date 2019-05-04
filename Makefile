URL=http://aphsai.github.io

default:
	@echo "Processing urls..."
	@go run hargonize.go 

fmt:
	@echo "\033[92mGo fmt source code...\033[0m"
	@go fmt *.go

