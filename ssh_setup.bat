winget install --id Cloudflare.cloudflared

echo "Host ssh.capstone.benmckallip.com\n   ProxyCommand C:\Program Files\Cloudflare\cloudflared.exe access ssh --hostname %h" >> C:\Users\%USERNAME%\.ssh\config