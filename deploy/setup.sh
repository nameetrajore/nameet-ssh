#!/bin/bash
# Run on the VM as root. Moves admin sshd to 2222, installs app on 22.
set -euo pipefail
mkdir -p /etc/ssh/sshd_config.d
echo "Port 2222" > /etc/ssh/sshd_config.d/99-port.conf
# Ubuntu 24.04 uses socket activation
if systemctl list-unit-files | grep -q '^ssh.socket'; then
  mkdir -p /etc/systemd/system/ssh.socket.d
  printf '[Socket]\nListenStream=\nListenStream=2222\n' > /etc/systemd/system/ssh.socket.d/port.conf
  systemctl daemon-reload; systemctl restart ssh.socket
fi
systemctl restart ssh || systemctl restart sshd
id nameetssh &>/dev/null || useradd --system --no-create-home --shell /usr/sbin/nologin nameetssh
mkdir -p /var/lib/nameet-ssh && chown nameetssh: /var/lib/nameet-ssh
install -m 755 /tmp/nameet-ssh-linux /usr/local/bin/nameet-ssh
install -m 644 /tmp/nameet-ssh.service /etc/systemd/system/nameet-ssh.service
systemctl daemon-reload && systemctl enable --now nameet-ssh && systemctl restart nameet-ssh
systemctl --no-pager status nameet-ssh | head -5
