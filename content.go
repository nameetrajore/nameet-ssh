package main

// Edit this file to change what visitors see.

const name = "Nameet Rajore"
const tagline = "Forward Deployed Engineer @ realfast · ex-Amazon Payments · Bengaluru"

var about = []string{
	"Hi, I'm Nameet.",
	"",
	"I build software that ships into real businesses — currently as a",
	"Forward Deployed Engineer at realfast. Before that, payments at Amazon.",
	"",
	"I like terminals, neovim, and tools that get out of the way.",
}

type project struct{ title, desc string }

var projects = []project{
	{"ssh.nameet.in", "this thing — a Go + Wish app you're looking at right now"},
	{"graphify", "turn any folder of code/docs into a queryable knowledge graph"},
}

type link struct{ label, url string }

var links = []link{
	{"blog", "https://nameet.in"},
	{"github", "https://github.com/nameetrajore"},
	{"x", "https://x.com/nameetrajore"},
}
