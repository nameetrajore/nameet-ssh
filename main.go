package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
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
	statsPath := flag.String("stats", "stats.json", "where visitor count and high score live")
	flag.Parse()

	st := loadStats(*statsPath)

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
				return newModel(st, st.visit()), nil
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
