build:
	env GOOS=linux go build -ldflags="-s -w" -o bootstrap main.go

deploy_prod: build
	serverless deploy --stage prod
