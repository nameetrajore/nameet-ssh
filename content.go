package main

// Edit this file to change what visitors see. Mirrors nameet.in.

const name = "Nameet Rajore"
const email = "nameet.rajore@gmail.com"

var home = []string{
	"# Hey, I'm Nameet.",
	"",
	"Forward Deployed Engineer at realfast. I build software across the stack, from backend systems and data pipelines to AI agents, and I see it through end to end: from understanding the client's problem to shipping something that moves their numbers.",
	"",
	"Before that I was on Amazon's payments infrastructure team building Tier-1 systems. I studied CS and Physics at BITS Pilani, helped start 2WheelR, and worked on ML research at Samsung R&D and LLM quantization at Kinara.ai.",
	"",
	"I write about systems, security, and how things work under the hood at nameet.in.",
	"",
	"→ ./projects · ./about · ./contact · ./play · ./shell",
	"",
	"Building something kick-ass? Say hi.",
}

var about = []string{
	"# Hello, I'm Nameet.",
	"",
	"I'm a Forward Deployed Engineer at realfast. I write software every day, working directly in client codebases across backend, data pipelines, and AI agents. I also own projects end to end: scoping the problem with the client, building the solution, and getting it into production.",
	"",
	"Before this, I was a backend engineer at Amazon working on Payments infrastructure, building Tier-1 systems that are robust and reliable.",
	"",
	"I graduated with a dual degree in Computer Science and Physics from BITS Pilani in 2024. While I was there, I also became a founding member of 2WheelR, leading its marketing and design team. That's also where I made my segue into software.",
	"",
	"I then went deep on the other end of the stack, researching unsupervised learning at Samsung R&D and optimizing LLM quantization techniques for edge devices at Kinara.ai. As is evident from my history, I love diving deep into unknown territories and finding my way out.",
	"",
	"I like to build, own, and see things through.",
	"",
	"If you're building something kick-ass, feel free to reach out!",
}

type project struct {
	title, desc, url string
	stack            []string
}

var projects = []project{
	{"Monke Interpreter", "A tree-walking interpreter for the Monke programming language, written in Go.",
		"https://github.com/nameetrajore/monke-interpreter", []string{"Go"}},
	{"2BRented", "A full-stack bike rental platform for browsing, filtering, and booking bikes across major Indian cities, with separate portals for customers, owners, and admins.",
		"https://github.com/nameetrajore/2BRented", []string{"React", "Express.js", "PostgreSQL", "AWS S3", "Razorpay"}},
	{"Lipschitz Adaptive LR", "Experimental framework implementing a Lipschitz Adaptive Learning Rate (LALR) policy for training regression models, with empirical validation notebooks.",
		"https://github.com/nameetrajore/lipschitz-adaptive-lr", []string{"Python", "Jupyter", "NumPy", "scikit-learn"}},
	{"Path-Finding Visualiser", "Interactive visualiser for classic pathfinding and maze generation algorithms, with step-by-step animations and configurable grid settings.",
		"https://github.com/nameetrajore/path-finding-visualiser", []string{"React", "Material UI", "JavaScript"}},
	{"Hypercube Quicksort (OpenMP)", "Parallel in-memory implementation of the Hypercube Quicksort algorithm using OpenMP for high-performance distributed-style sorting on shared-memory systems.",
		"https://github.com/nameetrajore/hyperquicksort-openmp", []string{"C", "OpenMP"}},
	{"SignalAI", "An AI news curator that filters out hype, fear-mongering, and speculation from AI journalism to surface signal over noise.",
		"https://github.com/nameetrajore/signal-ai", []string{"JavaScript", "Node.js"}},
	{"ssh.nameet.in", "This thing. A terminal portfolio served over SSH.",
		"https://github.com/nameetrajore/nameet-ssh", []string{"Go", "Wish", "Bubble Tea"}},
}

type link struct{ label, value, url string }

var contacts = []link{
	{"email", email, "mailto:" + email},
	{"github", "github.com/nameetrajore", "https://github.com/nameetrajore"},
	{"x", "x.com/nameetrajore", "https://x.com/nameetrajore"},
	{"blog", "nameet.in", "https://nameet.in"},
}
