package main

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

type fakeShell struct {
	lines []string
	input string
	hist  []string
	hi    int
}

const prompt = "visitor@nameet.in:~$ "

func newShell() fakeShell {
	return fakeShell{lines: []string{
		"fakesh 1.0 — nothing here is real, everything here is fun.",
		"type `help` to see what you can do.",
		"",
	}}
}

func (m *model) shellKey(msg tea.KeyPressMsg) tea.Cmd {
	sh := &m.shell
	switch msg.String() {
	case "esc":
		m.goTo(pHome)
	case "tab":
		m.goTo(m.page + 1)
	case "shift+tab":
		m.goTo(m.page - 1)
	case "enter":
		in := strings.TrimSpace(sh.input)
		sh.lines = append(sh.lines, prompt+sh.input)
		sh.input = ""
		if in != "" {
			sh.hist = append(sh.hist, in)
		}
		sh.hi = len(sh.hist)
		cmd := m.run(in)
		if len(sh.lines) > 300 {
			sh.lines = sh.lines[len(sh.lines)-300:]
		}
		return cmd
	case "backspace":
		if r := []rune(sh.input); len(r) > 0 {
			sh.input = string(r[:len(r)-1])
		}
	case "ctrl+u":
		sh.input = ""
	case "ctrl+l":
		sh.lines = nil
	case "up":
		if sh.hi > 0 {
			sh.hi--
			sh.input = sh.hist[sh.hi]
		}
	case "down":
		if sh.hi < len(sh.hist)-1 {
			sh.hi++
			sh.input = sh.hist[sh.hi]
		} else {
			sh.hi, sh.input = len(sh.hist), ""
		}
	default:
		if msg.Text != "" && len([]rune(sh.input)) < 120 {
			sh.input += msg.Text
		}
	}
	return nil
}

func (m *model) out(s ...string) { m.shell.lines = append(m.shell.lines, s...) }

func (m *model) run(in string) tea.Cmd {
	f := strings.Fields(in)
	if len(f) == 0 {
		return nil
	}
	cmd, args := strings.ToLower(f[0]), f[1:]
	rest := strings.Join(args, " ")
	switch cmd {
	case "help", "man", "?":
		m.out(
			"  ls, cd <page>     look around / go somewhere",
			"  whoami, pwd, date the usual",
			"  neofetch          system info, obviously",
			"  hire              the best command",
			"  cat <file>        read a file",
			"  echo, clear       you know these",
			"  coffee, ping      try them",
			"  snake             play a game",
			"  exit              leave (sad)",
			"  ...and a few secrets.",
		)
	case "ls":
		if rest == "-a" || rest == "-la" || rest == "-al" {
			m.out(".  ..  .secret  home/  projects/  about/  contact/  play/  shell/  secret.txt")
		} else {
			m.out("home/  projects/  about/  contact/  play/  shell/  secret.txt")
		}
	case "cd":
		target := strings.Trim(strings.TrimPrefix(rest, "./"), "/")
		switch target {
		case "", "~":
			m.goTo(pHome)
		case "..", "/":
			m.out("you can't leave that easily. (ctrl+c works though)")
		default:
			for i, n := range pageNames {
				if n == target {
					m.goTo(page(i))
					return nil
				}
			}
			m.out("cd: no such directory: " + rest)
		}
	case "cat", "less", "more", "bat":
		switch strings.TrimPrefix(rest, "./") {
		case "secret.txt":
			m.out("the real secret is ↑ ↑ ↓ ↓ ← → ← → b a (not in here, out there)")
		case ".secret":
			m.out("you found the hidden file. you'd be a great hire. so would he. → `hire`")
		case "":
			m.out("cat: meow.")
		default:
			m.out("cat: " + rest + ": no such file. try `ls`")
		}
	case "whoami":
		m.out(fmt.Sprintf("visitor #%d. a person of taste, clearly.", m.visitor))
	case "pwd":
		m.out("/home/visitor")
	case "date":
		m.out(time.Now().In(ist).Format("Mon 2 Jan 2006, 15:04 IST") + " — " + mood(time.Now().In(ist).Hour()))
	case "echo":
		m.out(rest)
	case "clear", "cls":
		m.shell.lines = nil
	case "neofetch", "fastfetch":
		m.out(
			"   ▄▄▄▄▄     visitor@nameet.in",
			"  █ ◕ ◕ █    ─────────────────",
			"  █  ▿  █    OS: nameet.in (wish v2)",
			"   ▀▀▀▀▀     Host: e2-micro, us-central1",
			"             Uptime: "+uptime(),
			"             Shell: fakesh 1.0",
			"             Theme: #FFA86A on #222129",
			fmt.Sprintf("             Visitors: %d", m.st.visits()),
		)
	case "hire":
		m.out("excellent choice. → "+email, "(copied to your clipboard, if your terminal allows it)")
		return tea.Batch(tea.SetClipboard(email), m.setFlash("copied "+email))
	case "email", "contact", "mail":
		m.out(email)
	case "sudo", "su", "doas":
		m.out("visitor is not in the sudoers file. This incident will be reported.")
	case "rm":
		m.out("nice try.")
	case "vim", "vi", "nvim":
		m.out("you're trapped now. (jk — nameet uses neovim btw)")
	case "emacs", "nano", "code":
		m.out("this is a neovim household.")
	case "coffee", "brew":
		m.out("   ( (", "    ) )", "  ........", "  |      |]", "  \\      /", "   `----'   here you go ☕")
	case "ping":
		m.out("pong")
	case "snake", "play", "game":
		m.goTo(pPlay)
	case "exit", "logout", "quit", ":q", ":wq", ":q!":
		return tea.Quit
	case "hello", "hi", "hey":
		m.out("hi! 👋 try `hire`, `neofetch` or `snake`.")
	case "42":
		m.out("the answer, yes. but what was the question?")
	case "make":
		if rest == "coffee" {
			m.out("make: *** no rule to make target 'coffee'. try `coffee`.")
		} else {
			m.out("make: *** no targets. stop.")
		}
	case "git":
		m.out("source lives at github.com/nameetrajore/nameet-ssh")
	default:
		m.out("fakesh: command not found: " + cmd + ". try `help`")
	}
	return nil
}

func uptime() string {
	d := time.Since(startedAt)
	switch {
	case d < time.Minute:
		return "just woke up"
	case d < time.Hour:
		return fmt.Sprintf("%d mins", int(d.Minutes()))
	case d < 48*time.Hour:
		return fmt.Sprintf("%d hours, %d mins", int(d.Hours()), int(d.Minutes())%60)
	}
	return fmt.Sprintf("%d days", int(d.Hours()/24))
}
