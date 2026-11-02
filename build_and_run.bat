if not exist config.env (
    echo Missing config.env. Copy config.env.example to config.env and fill in real values.
    exit /b 1
)
for /f "usebackq eol=# tokens=1,* delims==" %%A in ("config.env") do set "%%A=%%B"
call npm run build
cd /d "server" && go mod tidy && go build -o server . && go run server.go