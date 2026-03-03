# Kill any existing Go server instances
pkill -f '\./server' || true

cd 'Svelte Examples/plain-svelte-app' && npm run build
cd '../../Server Examples' && go mod tidy && go build server.go && chmod +x server && ./server > /dev/null 2>&1 &