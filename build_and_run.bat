call npm run build
cd /d "server" && go mod tidy && go build -o server . && go run server.go