# Kill any existing Go server instances
sudo pkill -f '\./server' || true

# Required for Google OAuth login.
# Example:
export GOOGLE_CLIENT_ID='1079255593069-2q1tq5l5jjdphphqoi3ekr4kiu0vmcb6.apps.googleusercontent.com '
export GOOGLE_CLIENT_SECRET='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'
export GOOGLE_REDIRECT_URL='https://capstone.benmckallip.com/auth/callback'

cd 'Svelte Examples/plain-svelte-app' && npm run build
cd '../../Server Examples' && go mod tidy && go build server.go && chmod +x server && ./server > server.log 2>&1 &