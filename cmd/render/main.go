package main

import (
	"log"

	"github.com/operator-sdk-demo/slide-controller/pkg/mdparser"
	"github.com/operator-sdk-demo/slide-controller/pkg/mdrender"
)

// Slide represents a single slide in the presentation
func main() {

	// Example slides data
	presentation := mdrender.Presentation{
		Slides: []mdrender.Slide{
			{
				Title:   "Kubernetes re-cap",
				Bullets: []string{"First bullet point", "Second bullet point"},
				Image:   "https://iximiuz.com/kubernetes-operator-pattern/kube-control-loop-3000-opt.png",
			},
			{
				Title:   "Next page",
				Bullets: []string{"Another bullet point", "Final bullet point"},
			},
		},
	}

	render := mdrender.RenderMarkdown(&presentation)

	if err := mdparser.CreateMarkdownParser(render); err != nil {
		log.Fatalf("failed to generate markdownParser")
	}
}
