#!/usr/bin/env bash
set -euo pipefail

# Usage: build_and_run.sh [--build-only | --no-build] [--foreground]
#   (no flags)     build, then restart this checkout's server in the background (logs to server/server.log)
#   --build-only   build the frontend and server binary, but don't start anything
#   --no-build     skip the build and start the existing server binary
#   --foreground   run the server in the foreground (for supervisors like systemd) instead of backgrounding it
DO_BUILD=true
DO_RUN=true
FOREGROUND=false
for arg in "$@"; do
	case "$arg" in
		--build-only) DO_RUN=false ;;
		--no-build) DO_BUILD=false ;;
		--foreground) FOREGROUND=true ;;
		*) echo "Unknown option: $arg"; exit 1 ;;
	esac
done

# Physical path, so the instance check below matches however the folder was reached (e.g. ~/capstone symlinks).
ROOT_DIR="$(cd "$(dirname "$0")" && pwd -P)"
SERVER_DIR="$ROOT_DIR/server"
SERVER_BIN="$SERVER_DIR/server"
cd "$ROOT_DIR"

# Load config and secrets (port, OAuth, integration, database) from the gitignored config file.
CONFIG_FILE="$ROOT_DIR/config.env"
if [ ! -f "$CONFIG_FILE" ] && [ -f "$ROOT_DIR/credentials.env" ]; then
	echo "Warning: credentials.env has been renamed to config.env; please rename your copy."
	CONFIG_FILE="$ROOT_DIR/credentials.env"
fi
if [ ! -f "$CONFIG_FILE" ]; then
	echo "Missing $CONFIG_FILE. Copy config.env.example to config.env and fill in real values."
	exit 1
fi
while IFS='=' read -r key value || [ -n "$key" ]; do
	case "$key" in ''|\#*) continue ;; esac
	export "$key=$value"
done < "$CONFIG_FILE"

if $DO_BUILD; then
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

	# Fresh checkouts (e.g. release tarballs) don't include node_modules.
	if [ ! -d node_modules ]; then
		npm ci
	fi
	npm run build
	cp Webpages/api-test.html Webpages/dist/
	(cd "$SERVER_DIR" && go mod tidy && go build -o server . && chmod +x server)
fi

if ! $DO_RUN; then
	exit 0
fi

if [ ! -x "$SERVER_BIN" ]; then
	echo "No server binary at $SERVER_BIN. Run without --no-build first."
	exit 1
fi

# Stop only the server started from this checkout (matched by its absolute path), so other
# instances on this machine (e.g. production) are left alone. Use sudo interactively so a
# teammate can replace an instance started by someone else.
if [ -t 0 ]; then
	sudo pkill -x -f "$SERVER_BIN" || true
else
	pkill -x -f "$SERVER_BIN" || true
fi
for _ in $(seq 1 20); do
	pgrep -x -f "$SERVER_BIN" >/dev/null || break
	sleep 0.5
done

# The server resolves the repo root from its working directory, so it must start in server/.
cd "$SERVER_DIR"
if $FOREGROUND; then
	exec "$SERVER_BIN"
fi
"$SERVER_BIN" > server.log 2>&1 &
echo "Server started on port ${PORT:-8080} (pid $!), logging to $SERVER_DIR/server.log"
