#!/bin/bash
# Rebuild and push to ssh.nameet.in (Oracle Cloud, admin SSH on port 2222).
set -euo pipefail
cd "$(dirname "$0")"
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/nameet-ssh-linux .
H="ubuntu@ssh.nameet.in"
scp -P 2222 -o StrictHostKeyChecking=accept-new dist/nameet-ssh-linux $H:/tmp/
ssh -p 2222 $H "sudo install -m 755 /tmp/nameet-ssh-linux /usr/local/bin/nameet-ssh && sudo systemctl restart nameet-ssh && systemctl is-active nameet-ssh"
