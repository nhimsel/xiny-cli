package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/net/html"
)

const baseURL = "https://learnxinyminutes.com/"

func main() {
	len := len(os.Args)

	if len > 3 || len < 2 {
		log.Fatalf("usage: %s <lang> [<locale>]",
			os.Args[0])
	}

	url := ""
	if len == 3 {
		url = baseURL + os.Args[2] + "/" + os.Args[1]
	} else {
		url = baseURL + os.Args[1]
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("HTTP error: %s", resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var output strings.Builder
	extractText(doc, &output, false)

	pagerPath, err := exec.LookPath("less")
	if err != nil {
		os.Stdout.WriteString(output.String())
	} else {
		pager := exec.Command(pagerPath, "-R")
		pager.Stdin = strings.NewReader(output.String())
		pager.Stdout = os.Stdout
		pager.Stderr = os.Stderr
		if err := pager.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func extractText(n *html.Node, output *strings.Builder, inPre bool) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, attr := range n.Attr {
			if attr.Key == "class" && attr.Val == "lang-choice" {
				return
			}
		}
	}

	if n.Type == html.ElementNode {
		switch n.Data {
		case "script", "style", "noscript":
			return
		case "pre":
			inPre = true
			output.WriteString("\n")
		}
	}

	if n.Type == html.TextNode {
		if inPre {
			output.WriteString(n.Data)
		} else {
			text := strings.Join(strings.Fields(n.Data), " ")
			if text != "" {
				output.WriteString(text)
				output.WriteByte(' ')
			}
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		extractText(child, output, inPre)
	}

	if n.Type == html.ElementNode {
		switch n.Data {
		case "p", "h1", "h2", "h3", "h4", "h5", "h6",
			"li", "section", "div":
			if !inPre {
				output.WriteByte('\n')
			}
		case "pre":
			output.WriteString("\n")
		}
	}
}
