package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type point struct{ x, y int }

type snakeState int

const (
	snakeReady snakeState = iota
	snakePlaying
	snakePaused
	snakeDead
)

type snakeTickMsg struct{ gen int }

type snakeGame struct {
	cols, rows int
	body       []point // body[0] is the head
	dir, next  point
	food       point
	score      int
	state      snakeState
	gen        int
	newHigh    bool
	rng        *rand.Rand
}

func (g *snakeGame) reset(cols, rows int) {
	if g.rng == nil {
		g.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	g.cols, g.rows = cols, rows
	cx, cy := cols/2, rows/2
	g.body = []point{{cx, cy}, {cx - 1, cy}, {cx - 2, cy}}
	g.dir, g.next = point{1, 0}, point{1, 0}
	g.score, g.newHigh = 0, false
	g.placeFood()
}

func (g *snakeGame) placeFood() {
	for {
		p := point{g.rng.Intn(g.cols), g.rng.Intn(g.rows)}
		if !g.occupies(p) {
			g.food = p
			return
		}
	}
}

func (g *snakeGame) occupies(p point) bool {
	for _, b := range g.body {
		if b == p {
			return true
		}
	}
	return false
}

// steer handles a direction key; returns false if the key isn't a direction.
func (g *snakeGame) steer(k string) bool {
	var d point
	switch k {
	case "up", "w", "k":
		d = point{0, -1}
	case "down", "s", "j":
		d = point{0, 1}
	case "left", "a", "h":
		d = point{-1, 0}
	case "right", "d", "l":
		d = point{1, 0}
	default:
		return false
	}
	if d.x != -g.dir.x || d.y != -g.dir.y { // no reversing into yourself
		g.next = d
	}
	return true
}

// step advances one frame; returns true if the snake died.
func (g *snakeGame) step() bool {
	g.dir = g.next
	h := g.body[0]
	n := point{(h.x + g.dir.x + g.cols) % g.cols, (h.y + g.dir.y + g.rows) % g.rows} // wrap around walls
	if g.occupies(n) && n != g.body[len(g.body)-1] {
		return true
	}
	g.body = append([]point{n}, g.body...)
	if n == g.food {
		g.score++
		g.placeFood()
	} else {
		g.body = g.body[:len(g.body)-1]
	}
	return false
}

func (g *snakeGame) tick() tea.Cmd {
	d := 130*time.Millisecond - time.Duration(g.score)*4*time.Millisecond
	if d < 55*time.Millisecond {
		d = 55 * time.Millisecond
	}
	gen := g.gen
	return tea.Tick(d, func(time.Time) tea.Msg { return snakeTickMsg{gen} })
}

func (g *snakeGame) start() tea.Cmd {
	g.state = snakePlaying
	g.gen++
	return g.tick()
}

func (g *snakeGame) view(high, highBy int, party bool, hue int) string {
	var b strings.Builder
	cells := make(map[point]int, len(g.body))
	for i, p := range g.body {
		cells[p] = i + 1
	}
	head := fgStyle
	if party {
		head = lipgloss.NewStyle().Foreground(partyColors[hue%len(partyColors)])
	}
	for y := 0; y < g.rows; y++ {
		for x := 0; x < g.cols; x++ {
			p := point{x, y}
			switch {
			case cells[p] == 1:
				b.WriteString(head.Render("██"))
			case cells[p] > 1:
				b.WriteString(visitedStyle.Render("▓▓"))
			case p == g.food:
				b.WriteString(foodStyle.Render(" ◆"))
			default:
				b.WriteString(faintStyle.Render(" ·"))
			}
		}
		if y < g.rows-1 {
			b.WriteByte('\n')
		}
	}
	field := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(dimColor).Render(b.String())

	hs := "none yet — be the first"
	if high > 0 {
		hs = fmt.Sprintf("%d by visitor #%d", high, highBy)
	}
	top := fgStyle.Render(fmt.Sprintf("score %d", g.score)) + dimStyle.Render("  ·  all-time high "+hs)

	var msg string
	switch g.state {
	case snakeReady:
		msg = "space to start · eat the ◆ · walls wrap"
	case snakePaused:
		msg = "paused · space to resume"
	case snakeDead:
		msg = fmt.Sprintf("game over · score %d · space to retry", g.score)
		if g.newHigh {
			msg = "★ new all-time high score! ★ · space to retry"
		}
	case snakePlaying:
		msg = " "
	}
	return lipgloss.JoinVertical(lipgloss.Left, "# Snake", "", top, field, fgStyle.Render(msg))
}
