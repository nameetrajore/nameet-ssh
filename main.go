package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"charm.land/ssh"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	bm "charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	gossh "golang.org/x/crypto/ssh"
)

func main() {
	addr := flag.String("addr", ":2222", "listen address")
	hostKey := flag.String("hostkey", ".ssh/id_ed25519", "host key path (created if missing)")
	flag.Parse()

	s, err := wish.NewServer(
		wish.WithAddress(*addr),
		wish.WithHostKeyPath(*hostKey),
		// Anyone may connect: accept any key and keyboard-interactive (no password prompt).
		wish.WithPublicKeyAuth(func(ssh.Context, ssh.PublicKey) bool { return true }),
		wish.WithKeyboardInteractiveAuth(func(ssh.Context, gossh.KeyboardInteractiveChallenge) bool { return true }),
		wish.WithIdleTimeout(10*time.Minute),
		wish.WithMaxTimeout(30*time.Minute),
		wish.WithMiddleware(
			bm.Middleware(func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
				return model{}, nil
			}),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatal("server", "err", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	log.Info("listening", "addr", *addr)
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("serve", "err", err)
			done <- nil
		}
	}()
	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = s.Shutdown(ctx)
}

// ---- TUI ----

var tabs = []string{"about", "projects", "links"}

type model struct {
	tab, w, h int
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "right", "l", "tab":
			m.tab = (m.tab + 1) % len(tabs)
		case "left", "h", "shift+tab":
			m.tab = (m.tab + len(tabs) - 1) % len(tabs)
		case "1", "2", "3":
			m.tab = int(msg.String()[0] - '1')
		}
	}
	return m, nil
}

var (
	accent = lipgloss.Color("#F25D94")
	muted  = lipgloss.Color("#7D7D8A")
	title  = lipgloss.NewStyle().Bold(true).Foreground(accent)
	dim    = lipgloss.NewStyle().Foreground(muted)
	active = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#1A1A1A")).Background(accent).Padding(0, 1)
	idle   = lipgloss.NewStyle().Foreground(muted).Padding(0, 1)
	box    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(accent).Padding(1, 3)
)

func (m model) View() tea.View {
	var b strings.Builder
	b.WriteString(title.Render(name) + "\n")
	b.WriteString(dim.Render(tagline) + "\n\n")

	var t []string
	for i, n := range tabs {
		if i == m.tab {
			t = append(t, active.Render(n))
		} else {
			t = append(t, idle.Render(n))
		}
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, t...) + "\n\n")

	switch m.tab {
	case 0:
		b.WriteString(strings.Join(about, "\n"))
	case 1:
		for _, p := range projects {
			b.WriteString(lipgloss.NewStyle().Bold(true).Render(p.title) + "\n  " + dim.Render(p.desc) + "\n\n")
		}
	case 2:
		for _, l := range links {
			b.WriteString(lipgloss.NewStyle().Width(8).Bold(true).Render(l.label) + lipgloss.NewStyle().Foreground(accent).Render(l.url) + "\n")
		}
	}
	b.WriteString("\n\n" + dim.Render("←/→ switch · q quit"))

	bw := 82
	if m.w > 0 && m.w-2 < bw {
		bw = m.w - 2
	}
	content := box.Width(bw).Render(strings.TrimRight(b.String(), "\n"))
	v := tea.NewView(lipgloss.Place(m.w, m.h, lipgloss.Center, lipgloss.Center, content))
	v.AltScreen = true
	return v
}
