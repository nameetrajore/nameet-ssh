#!/bin/bash
# First-time setup on an Ubuntu 24.04 VM (run as root, with nameet-ssh-linux and
# nameet-ssh.service in /tmp). Moves admin sshd to 2222, runs the app on 22.
# Verify you can reach port 2222 (cloud firewall + host firewall) BEFORE relying on it.
set -euo pipefail
# Oracle images ship iptables with a trailing REJECT; open 2222 above it.
if iptables -L INPUT -n | grep -q REJECT; then
  pos=$(iptables -L INPUT -n --line-numbers | awk '/REJECT/{print $1; exit}')
  iptables -I INPUT "$pos" -p tcp -m state --state NEW --dport 2222 -j ACCEPT
  netfilter-persistent save || true
fi
echo "Port 2222" > /etc/ssh/sshd_config.d/99-port.conf
systemctl disable --now ssh.socket 2>/dev/null || true
rm -rf /etc/systemd/system/ssh.socket.d
systemctl daemon-reload
systemctl enable ssh.service
systemctl restart ssh.service
id nameetssh &>/dev/null || useradd --system --no-create-home --shell /usr/sbin/nologin nameetssh
mkdir -p /var/lib/nameet-ssh && chown nameetssh: /var/lib/nameet-ssh
install -m 755 /tmp/nameet-ssh-linux /usr/local/bin/nameet-ssh
install -m 644 /tmp/nameet-ssh.service /etc/systemd/system/nameet-ssh.service
systemctl daemon-reload && systemctl enable --now nameet-ssh
systemctl --no-pager status nameet-ssh | head -5
