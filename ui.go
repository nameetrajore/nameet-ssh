package main

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ---- theme: matches nameet.in ----

var (
	bgColor      = lipgloss.Color("#222129")
	fgColor      = lipgloss.Color("#FFA86A")
	visitedColor = lipgloss.Color("#bc8d6b")
	dimColor     = lipgloss.Color("#A77250") // fg at ~60% over bg, like the site's footer
	faintColor   = lipgloss.Color("#4a3a38")

	fgStyle      = lipgloss.NewStyle().Foreground(fgColor)
	boldStyle    = lipgloss.NewStyle().Foreground(fgColor).Bold(true)
	dimStyle     = lipgloss.NewStyle().Foreground(dimColor)
	faintStyle   = lipgloss.NewStyle().Foreground(faintColor)
	visitedStyle = lipgloss.NewStyle().Foreground(visitedColor)
	foodStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFE0C2")).Bold(true)
	invStyle     = lipgloss.NewStyle().Foreground(bgColor).Background(fgColor)
	chipStyle    = lipgloss.NewStyle().Foreground(bgColor).Background(visitedColor).Padding(0, 1)
	boxStyle     = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(fgColor).Padding(1, 3)

	partyColors = lipgloss.Blend1D(16,
		lipgloss.Color("#FFA86A"), lipgloss.Color("#FF6AC1"), lipgloss.Color("#6AE0FF"),
		lipgloss.Color("#B4FF6A"), lipgloss.Color("#FFA86A"))

	ist = time.FixedZone("IST", 5*3600+1800)
)

type page int

const (
	pHome page = iota
	pProjects
	pAbout
	pContact
	pPlay
	pShell
)

var pageNames = []string{"home", "projects", "about", "contact", "play", "shell"}

var konami = []string{"up", "up", "down", "down", "left", "right", "left", "right", "b", "a"}

type (
	blinkMsg      struct{}
	bootTickMsg   struct{}
	bootEndMsg    struct{}
	flashClearMsg struct{ id int }
)

type model struct {
	w, h    int
	st      *stats
	visitor int
	page    page

	booting bool
	bootN   int
	boot    []rune

	blink   bool
	help    bool
	scroll  int
	projSel int
	conSel  int
	flash   string
	flashID int
	party   bool
	hue     int
	keys    []string

	snake snakeGame
	shell fakeShell
}

func newModel(st *stats, visitor int) model {
	boot := fmt.Sprintf("$ ssh ssh.nameet.in\n"+
		"  resolving nameet.in .............. ok\n"+
		"  authenticating ................... no password needed, you're cool\n"+
		"  brewing coffee ................... ok\n"+
		"  welcome, visitor #%d.", visitor)
	return model{st: st, visitor: visitor, booting: true, boot: []rune(boot), shell: newShell()}
}

func blinkTick() tea.Cmd {
	return tea.Tick(530*time.Millisecond, func(time.Time) tea.Msg { return blinkMsg{} })
}

func bootTick() tea.Cmd {
	return tea.Tick(14*time.Millisecond, func(time.Time) tea.Msg { return bootTickMsg{} })
}

func (m *model) setFlash(s string) tea.Cmd {
	m.flashID++
	m.flash = s
	id := m.flashID
	return tea.Tick(2500*time.Millisecond, func(time.Time) tea.Msg { return flashClearMsg{id} })
}

func (m model) Init() tea.Cmd { return tea.Batch(bootTick(), blinkTick()) }

func (m *model) goTo(p page) {
	n := page(len(pageNames))
	p = (p%n + n) % n
	if m.page == pPlay && p != pPlay && m.snake.state == snakePlaying {
		m.snake.state = snakePaused
	}
	m.page, m.scroll, m.help = p, 0, false
	if p == pPlay && m.snake.state != snakePlaying && m.snake.state != snakePaused {
		if cols, rows := m.fieldSize(); cols >= 8 && rows >= 5 {
			m.snake.reset(cols, rows)
			m.snake.state = snakeReady
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
	case blinkMsg:
		m.blink = !m.blink
		m.hue++
		return m, blinkTick()
	case bootTickMsg:
		if !m.booting {
			return m, nil
		}
		m.bootN += 3
		if m.bootN >= len(m.boot) {
			m.bootN = len(m.boot)
			return m, tea.Tick(900*time.Millisecond, func(time.Time) tea.Msg { return bootEndMsg{} })
		}
		return m, bootTick()
	case bootEndMsg:
		m.booting = false
	case flashClearMsg:
		if msg.id == m.flashID {
			m.flash = ""
		}
	case snakeTickMsg:
		g := &m.snake
		if msg.gen != g.gen || g.state != snakePlaying {
			return m, nil
		}
		if g.step() {
			g.state = snakeDead
			g.newHigh = m.st.submit(g.score, m.visitor)
			return m, nil
		}
		return m, g.tick()
	case tea.MouseClickMsg:
		return m.click(msg.Mouse())
	case tea.MouseWheelMsg:
		switch msg.Mouse().Button {
		case tea.MouseWheelUp:
			m.move(-1)
		case tea.MouseWheelDown:
			m.move(1)
		}
	case tea.KeyPressMsg:
		return m.key(msg)
	}
	return m, nil
}

func (m model) key(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	k := msg.String()
	if k == "ctrl+c" {
		return m, tea.Quit
	}
	if m.booting {
		m.booting = false
		return m, nil
	}

	var cmds []tea.Cmd
	m.keys = append(m.keys, k)
	if len(m.keys) > len(konami) {
		m.keys = m.keys[len(m.keys)-len(konami):]
	}
	if strings.Join(m.keys, " ") == strings.Join(konami, " ") {
		m.party = !m.party
		m.keys = nil
		if m.party {
			cmds = append(cmds, m.setFlash("★ party mode unlocked ★"))
		} else {
			cmds = append(cmds, m.setFlash("party's over"))
		}
	}

	if m.help {
		m.help = false
		return m, tea.Batch(cmds...)
	}
	if m.page == pShell {
		cmds = append(cmds, m.shellKey(msg))
		return m, tea.Batch(cmds...)
	}
	if m.page == pPlay && m.snake.state == snakePlaying {
		if m.snake.steer(k) {
			return m, tea.Batch(cmds...)
		}
	}

	switch k {
	case "q":
		return m, tea.Quit
	case "?":
		m.help = true
	case "right", "l", "tab":
		m.goTo(m.page + 1)
	case "left", "h", "shift+tab":
		m.goTo(m.page - 1)
	case "1", "2", "3", "4", "5", "6":
		m.goTo(page(k[0] - '1'))
	case "up", "k":
		m.move(-1)
	case "down", "j":
		m.move(1)
	case "enter", "space", "o", "y":
		cmds = append(cmds, m.activate(k))
	case "p":
		if m.page == pPlay && m.snake.state == snakePlaying {
			m.snake.state = snakePaused
		}
	case "esc":
		if m.page == pPlay && m.snake.state == snakePlaying {
			m.snake.state = snakePaused
		} else {
			m.goTo(pHome)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *model) move(d int) {
	switch m.page {
	case pProjects:
		m.projSel = clamp(m.projSel+d, 0, len(projects)-1)
	case pContact:
		m.conSel = clamp(m.conSel+d, 0, len(contacts)-1)
	case pShell, pPlay:
	default:
		m.scroll = max(0, m.scroll+d)
	}
}

func (m *model) activate(k string) tea.Cmd {
	switch m.page {
	case pProjects:
		u := projects[m.projSel].url
		return tea.Batch(tea.SetClipboard(u), m.setFlash("copied "+u))
	case pContact:
		c := contacts[m.conSel]
		return tea.Batch(tea.SetClipboard(c.value), m.setFlash("copied "+c.value))
	case pPlay:
		if k != "space" && k != "enter" {
			return nil
		}
		g := &m.snake
		switch g.state {
		case snakeReady, snakeDead:
			cols, rows := m.fieldSize()
			if cols < 8 || rows < 5 {
				return m.setFlash("terminal too small for snake — make it bigger")
			}
			g.reset(cols, rows)
			return g.start()
		case snakePaused:
			return g.start()
		}
	}
	return nil
}

// ---- layout ----

func (m model) frameWidth() int { return min(m.w-2, 86) }
func (m model) innerWidth() int { return m.frameWidth() - 8 } // border 2 + padding 6
func (m model) bodyHeight() int { return max(4, m.h-12) }     // everything that isn't the body
func (m model) leftMargin() int { return max(0, (m.w-m.frameWidth())/2) }
func (m model) fieldSize() (int, int) {
	return min(30, (m.innerWidth()-2)/2), min(16, m.bodyHeight()-5)
}

const navRow = 3 // margin, title, blank, nav

func (m model) navSpans() [][2]int {
	var spans [][2]int
	x := 0
	for _, n := range pageNames {
		w := len("./"+n) + 2
		spans = append(spans, [2]int{x, x + w})
		x += w + 1
	}
	return spans
}

func (m model) click(mo tea.Mouse) (tea.Model, tea.Cmd) {
	if m.booting {
		m.booting = false
		return m, nil
	}
	if mo.Y == navRow {
		x := mo.X - m.leftMargin()
		for i, s := range m.navSpans() {
			if x >= s[0] && x < s[1] {
				m.goTo(page(i))
			}
		}
	}
	return m, nil
}

func (m model) titleColor() lipgloss.Style {
	if m.party {
		return lipgloss.NewStyle().Foreground(bgColor).Background(partyColors[m.hue%len(partyColors)])
	}
	return invStyle
}

func (m model) nav() string {
	var parts []string
	for i, n := range pageNames {
		if page(i) == m.page {
			parts = append(parts, m.titleColor().Padding(0, 1).Render("./"+n))
		} else {
			parts = append(parts, " "+dimStyle.Render("./")+fgStyle.Render(n)+" ")
		}
	}
	return strings.Join(parts, " ")
}

func wrap(w int, lines ...string) []string {
	var out []string
	st := lipgloss.NewStyle().Width(w)
	for _, l := range lines {
		if l == "" {
			out = append(out, "")
			continue
		}
		if strings.HasPrefix(l, "# ") {
			out = append(out, dimStyle.Render("# ")+boldStyle.Render(l[2:]))
			continue
		}
		out = append(out, strings.Split(st.Render(fgStyle.Render(l)), "\n")...)
	}
	return out
}

func (m model) body() ([]string, bool) {
	iw := m.innerWidth()
	if m.help {
		return wrap(iw,
			"# Keys",
			"",
			"←/→  h/l  tab       switch page",
			"1–6                 jump to page",
			"↑/↓  j/k  wheel     move / scroll",
			"enter               copy link / start game",
			"click               the ./nav works with a mouse",
			"esc                 back home",
			"q  ctrl+c           quit",
			"",
			"psst: there's a secret code. the shell might know it.",
		), false
	}
	switch m.page {
	case pHome:
		return wrap(iw, home...), true
	case pAbout:
		return wrap(iw, about...), true
	case pProjects:
		lines := []string{dimStyle.Render("# ") + boldStyle.Render("Projects"), ""}
		for i, p := range projects {
			if i != m.projSel {
				lines = append(lines, "  "+fgStyle.Render(p.title))
				continue
			}
			lines = append(lines, m.titleColor().Render("▸ "+p.title+" "))
			d := lipgloss.NewStyle().Width(iw - 4).Render(fgStyle.Render(p.desc))
			for _, l := range strings.Split(d, "\n") {
				lines = append(lines, "    "+l)
			}
			var chips []string
			for _, s := range p.stack {
				chips = append(chips, chipStyle.Render(s))
			}
			lines = append(lines, "    "+strings.Join(chips, " "))
			lines = append(lines, "    "+dimStyle.Hyperlink(p.url).Render("↗ "+strings.TrimPrefix(p.url, "https://")), "")
		}
		return lines, false
	case pContact:
		lines := []string{dimStyle.Render("# ") + boldStyle.Render("Say hi"), "",
			fgStyle.Render("Building something kick-ass? Reach out."), ""}
		for i, c := range contacts {
			label := lipgloss.NewStyle().Width(8).Render(c.label)
			if i == m.conSel {
				lines = append(lines, m.titleColor().Render("▸ "+label+c.value+" "))
			} else {
				lines = append(lines, "  "+dimStyle.Render(label)+fgStyle.Hyperlink(c.url).Render(c.value))
			}
		}
		lines = append(lines, "", dimStyle.Render("enter copies · or type `hire` in ./shell"))
		return lines, false
	case pPlay:
		cols, rows := m.fieldSize()
		if m.snake.cols == 0 || cols < 8 || rows < 5 {
			return wrap(iw, "# Snake", "", "your terminal is too small to play. make it bigger!"), false
		}
		hi, by := m.st.high()
		return strings.Split(m.snake.view(hi, by, m.party, m.hue), "\n"), false
	case pShell:
		var lines []string
		for _, l := range m.shell.lines {
			lines = append(lines, wrap(iw, l)...)
		}
		cur := " "
		if m.blink {
			cur = "█"
		}
		in := boldStyle.Render("visitor@nameet.in") + dimStyle.Render(":~$ ") + fgStyle.Render(m.shell.input+cur)
		lines = append(lines, in)
		if h := m.bodyHeight(); len(lines) > h {
			lines = lines[len(lines)-h:]
		}
		return lines, false
	}
	return nil, false
}

func mood(h int) string {
	switch {
	case h < 7:
		return "he's probably asleep"
	case h < 10:
		return "he's probably making coffee"
	case h < 19:
		return "he's probably shipping something"
	default:
		return "he's probably debugging something"
	}
}

func (m model) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	v.BackgroundColor = bgColor
	v.ForegroundColor = fgColor
	if m.w == 0 {
		return v
	}
	if m.w < 50 || m.h < 18 {
		v.SetContent(fgStyle.Render("\n  nameet.in\n\n  make your terminal a\n  little bigger :)\n\n  (needs 50×18)"))
		return v
	}

	cur := " "
	if m.blink {
		cur = "_"
	}

	if m.booting {
		s := string(m.boot[:m.bootN])
		var b strings.Builder
		for i, l := range strings.Split(s, "\n") {
			if i == 0 {
				b.WriteString(dimStyle.Render("$ ") + boldStyle.Render(strings.TrimPrefix(l, "$ ")))
			} else {
				b.WriteString(fgStyle.Render(l))
			}
			b.WriteByte('\n')
		}
		out := strings.TrimRight(b.String(), "\n") + fgStyle.Render("█")
		v.SetContent(lipgloss.NewStyle().Margin(1, 0, 0, m.leftMargin()).Render(out))
		return v
	}

	title := m.titleColor().Padding(0, 1).Render(name + cur)

	lines, scrollable := m.body()
	h := m.bodyHeight()
	more := ""
	if scrollable && len(lines) > h {
		maxScroll := len(lines) - h
		sc := min(m.scroll, maxScroll)
		lines = lines[sc : sc+h]
		switch {
		case sc < maxScroll && sc > 0:
			more = " · ↑↓ more"
		case sc < maxScroll:
			more = " · ↓ more"
		default:
			more = " · ↑ more"
		}
	}
	box := boxStyle.Width(m.frameWidth()).Render(strings.Join(lines, "\n"))

	hint := "←/→ pages · ↑/↓ move · enter select · ? help · q quit"
	switch m.page {
	case pShell:
		hint = "type a command · tab next page · esc home · ctrl+c quit"
	case pPlay:
		hint = "arrows/wasd steer · space start · p pause · tab next page · q quit"
	}
	if m.flash != "" {
		hint = fgStyle.Render("✓ " + m.flash)
	} else {
		hint = dimStyle.Render(hint + more)
	}
	now := time.Now().In(ist)
	status := dimStyle.Render(fmt.Sprintf("visitor #%d · %s in Bengaluru, %s", m.visitor, now.Format("15:04"), mood(now.Hour())))

	col := lipgloss.JoinVertical(lipgloss.Left, title, "", m.nav(), "", box, " "+hint, " "+status)
	v.SetContent(lipgloss.NewStyle().Margin(1, 0, 0, m.leftMargin()).Render(col))
	return v
}

func clamp(v, lo, hi int) int { return max(lo, min(hi, v)) }
