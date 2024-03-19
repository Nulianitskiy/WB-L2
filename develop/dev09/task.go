package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

type Wget struct {
	URL             string
	Depth           int
	OnlyInternal    bool
	Links           []string
	Path            string
	LinksVisited    map[string]struct{}
	FilesDownloaded int
}

func Start(args []string) error {
	wg := Wget{URL: args[len(args)-1], Links: []string{args[len(args)-1]}, LinksVisited: make(map[string]struct{}, 1000)}
	wg.LinksVisited[wg.URL] = struct{}{}
	fs := flag.NewFlagSet("wgetflags", flag.ContinueOnError)
	fs.IntVar(&wg.Depth, "r", -1, "depth")
	fs.BoolVar(&wg.OnlyInternal, "i", false, "download only internal pages")
	fs.StringVar(&wg.Path, "f", "", "where to store sites")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	parsedURL, err := url.Parse(wg.URL)
	if err != nil {
		return err
	}
	host := parsedURL.Hostname()

	wg.Path, err = wg.MkDir(wg.Path, host)
	if err != nil {
		return err
	}

	err = wg.Process()
	if err != nil {
		os.RemoveAll(wg.Path)
		return err
	}
	fmt.Println(wg.FilesDownloaded)

	return nil
}

func (wg *Wget) MkDir(parentDir string, newDir string) (string, error) {

	fullPath := filepath.Join(parentDir, newDir)

	err := os.Mkdir(fullPath, 0755)
	if err != nil {
		return "", err
	}
	return fullPath, nil
}

func (wg *Wget) DownloadSite(link string) ([]byte, error) {
	response, err := http.Get(wg.URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (wg *Wget) ParseHTML(urlS string, data []byte) ([]string, error) {
	page := bytes.NewReader(data)
	doc, err := html.Parse(page)
	if err != nil {
		return nil, err
	}

	pageURL, err := url.Parse(urlS)
	if err != nil {
		return nil, err
	}

	links := make([]string, 0, len(wg.Links)/4)
	extractLinks := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					linkURL, err := url.Parse(attr.Val)
					if err != nil {
						continue
					}

					if !linkURL.IsAbs() {
						linkURL = pageURL.ResolveReference(linkURL)
					}

					if linkURL.Host != "" && linkURL.Host != pageURL.Host && wg.OnlyInternal {
						continue
					}

					if _, ok := wg.LinksVisited[strings.TrimSuffix(linkURL.String(), "/")]; ok {
						continue
					} else {
						wg.LinksVisited[strings.TrimSuffix(linkURL.String(), "/")] = struct{}{}
					}

					links = append(links, linkURL.String())
				}
			}
		}
	}

	var stack []*html.Node
	stack = append(stack, doc)
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		extractLinks(node)
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			stack = append(stack, c)
		}
	}

	return links, nil
}

func (wg *Wget) Process() error {
	var needIncr bool
	if wg.Depth == -1 {
		wg.Depth = 0
		needIncr = true
	}
	for i := 0; i <= wg.Depth; i++ {
		if len(wg.Links) == 0 {
			return nil
		}
		nextDepthLinks := make([]string, 0, len(wg.Links))
		fmt.Println(i)

		dirName := fmt.Sprintf("level%d", i)
		dirName, err := wg.MkDir(wg.Path, dirName)
		if err != nil {
			return err
		}

		for _, v := range wg.Links {
			fmt.Println(v)
			htmlData, err := wg.DownloadSite(v)
			if err != nil {
				return err
			}

			err = wg.Save(htmlData, dirName, v)
			if err != nil {
				return err
			}

			links, err := wg.ParseHTML(v, htmlData)
			if err != nil {
				return err
			}
			nextDepthLinks = append(nextDepthLinks, links...)
			wg.FilesDownloaded++
			wg.LinksVisited[v] = struct{}{}
		}
		wg.Links = nextDepthLinks
		if needIncr {
			wg.Depth++
		}
	}
	return nil
}

func (wg *Wget) Save(data []byte, dir string, urlS string) error {
	urlS = strings.TrimPrefix(urlS, "https://")
	splitted := strings.Split(urlS, "/")
	urlS = strings.Join(splitted, "_")
	filePath := fmt.Sprintf("%s/%s.html", dir, urlS)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := Start(os.Args[1:])
	if err != nil {
		fmt.Println(err.Error())
	}
}
