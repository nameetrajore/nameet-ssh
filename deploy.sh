#!/bin/bash
# Rebuild and push to ssh.nameet.in (admin SSH is on port 2222).
set -euo pipefail
cd "$(dirname "$0")"
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/nameet-ssh-linux .
H="-i $HOME/.ssh/google_compute_engine -o StrictHostKeyChecking=accept-new -p 2222"
scp ${H/-p/-P} dist/nameet-ssh-linux nameetrajore@ssh.nameet.in:/tmp/
ssh $H nameetrajore@ssh.nameet.in "sudo install -m 755 /tmp/nameet-ssh-linux /usr/local/bin/nameet-ssh && sudo systemctl restart nameet-ssh && systemctl is-active nameet-ssh"
