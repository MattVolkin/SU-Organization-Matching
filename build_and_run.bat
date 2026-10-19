if not exist credentials.env (
    echo Missing credentials.env. Copy credentials.env.example to credentials.env and fill in real values.
    exit /b 1
)
for /f "usebackq eol=# tokens=1,* delims==" %%A in ("credentials.env") do set "%%A=%%B"
call npm run build
cd /d "server" && go mod tidy && go build -o server . && go run server.go