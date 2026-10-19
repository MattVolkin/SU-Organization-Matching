#!/usr/bin/env bash
set -euo pipefail

# Kill any existing Go server instances
sudo pkill -f '\./server' || true

# Load secrets (OAuth, integration, database) from the gitignored credentials file.
CREDENTIALS_FILE="$(dirname "$0")/credentials.env"
if [ ! -f "$CREDENTIALS_FILE" ]; then
	echo "Missing $CREDENTIALS_FILE. Copy credentials.env.example to credentials.env and fill in real values."
	exit 1
fi
while IFS='=' read -r key value || [ -n "$key" ]; do
	case "$key" in ''|\#*) continue ;; esac
	export "$key=$value"
done < "$CREDENTIALS_FILE"

# Load nvm when available so this script uses the project Node version.
export NVM_DIR="$HOME/.nvm"
if [ -s "$NVM_DIR/nvm.sh" ]; then
	. "$NVM_DIR/nvm.sh"
fi

if command -v nvm >/dev/null 2>&1; then
	nvm use 24 >/dev/null 2>&1 || true
fi

if ! command -v node >/dev/null 2>&1; then
	echo "Node.js is not installed. Install nvm and run: nvm install 24 && nvm use 24"
	exit 1
fi

NODE_MAJOR="$(node -p "process.versions.node.split('.')[0]")"
if [ "$NODE_MAJOR" -lt 24 ]; then
	echo "Node.js 24+ is required. Current: $(node -v)"
	echo "Update with: nvm install 24 && nvm alias default 24 && nvm use 24"
	exit 1
fi

npm run build
cp Webpages/api-test.html Webpages/dist/
cd 'server' && go mod tidy && go build -o server . && chmod +x server && ./server > server.log 2>&1 &