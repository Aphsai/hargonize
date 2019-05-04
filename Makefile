URL=aphsai.github.io

default:
	@echo "\033[92mProcessing ${URL}...\033[0m"
	@go run hargonize.go -url=${URL}

fmt:
	@echo "\033[92mGo fmt source code...\033[0m"
	@go fmt *.go

