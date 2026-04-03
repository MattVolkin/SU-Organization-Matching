# Add cloudflare gpg key
sudo mkdir -p --mode=0755 /usr/share/keyrings
curl -fsSL https://pkg.cloudflare.com/cloudflare-main.gpg | sudo tee /usr/share/keyrings/cloudflare-main.gpg >/dev/null

# Add this repo to your apt repositories (Stable only)
echo 'deb [signed-by=/usr/share/keyrings/cloudflare-main.gpg] https://pkg.cloudflare.com/cloudflared any main' | sudo tee /etc/apt/sources.list.d/cloudflared.list

# install cloudflared
sudo apt-get update && sudo apt-get install cloudflared

# Configure SSH config for Cloudflare Access
echo -e "\n\nHost capstone\n   HostName capstonessh.benmckallip.com\n   ProxyCommand /usr/bin/cloudflared access ssh --hostname %h" >> ~/.ssh/config

cloudflared access login capstonessh.benmckallip.com