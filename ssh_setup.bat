winget install --id Cloudflare.cloudflared

echo "Host capstone\n   HostName capstonessh.benmckallip.com\n   ProxyCommand C:\Program Files\Cloudflare\cloudflared.exe access ssh --hostname %h" >> C:\Users\%USERNAME%\.ssh\config

"C:\Program Files\Cloudflare\cloudflared.exe" access login capstonessh.benmckallip.com