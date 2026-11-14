# ssh.nameet.in

```
ssh ssh.nameet.in
```

A terminal portfolio served over SSH. Go + [Wish](https://github.com/charmbracelet/wish) + [Bubble Tea](https://github.com/charmbracelet/bubbletea). No shell — just the TUI.

- `content.go` — what visitors see
- `main.go` — SSH server + TUI
- `deploy/` — systemd unit + first-time VM setup (moves admin sshd to 2222, app to 22)
- `deploy.sh` — rebuild and push

Run locally: `go run . -addr :2222` then `ssh -p 2222 localhost`.
