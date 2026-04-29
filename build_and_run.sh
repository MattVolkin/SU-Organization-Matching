#!/usr/bin/env bash
set -euo pipefail

# Kill any existing Go server instances
sudo pkill -f '\./server' || true

# Required for Google OAuth login.
# Example:
export GOOGLE_CLIENT_ID='1079255593069-2q1tq5l5jjdphphqoi3ekr4kiu0vmcb6.apps.googleusercontent.com '
export GOOGLE_CLIENT_SECRET='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'
export GOOGLE_REDIRECT_URL='https://capstone.benmckallip.com/auth/callback'
export GOOGLE_INTEGRATION_TRIGGER_URL='https://integrations.googleapis.com/v2/projects/capstone-489704/locations/us-south1/integrations/email-on-pubsub-event:execute?triggerId=api_trigger/sendEmail'

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