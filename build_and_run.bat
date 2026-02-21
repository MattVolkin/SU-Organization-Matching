cd /d "Svelte Examples\plain-svelte-app" && call npm run build
cd /d "..\..\Server Examples" && go mod tidy && go run server.go