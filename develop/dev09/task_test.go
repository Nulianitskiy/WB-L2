package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var page1 = []byte(`<!-- index.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Index Page</title>
</head>
<body>
    <h1>Index Page</h1>
    <ul>
        <li><a href="page1.html">Page 1</a></li>
        <li><a href="page2.html">Page 2</a></li>
        <li><a href="https://example.com">Example Site</a></li>
        <li><a href="https://www.google.com">Google</a></li>
    </ul>
</body>
</html>`)
var page2 = []byte(`<!-- page1.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Page 1</title>
</head>
<body>
    <h1>Page 1</h1>
    <ul>
        <li><a href="index.html">Index Page</a></li>
        <li><a href="page2.html">Page 2</a></li>
        <li><a href="https://example.com">Example Site</a></li>
        <li><a href="https://www.google.com">Google</a></li>
    </ul>
</body>
</html>`)

var page3 = []byte(`<!-- page2.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Page 2</title>
</head>
<body>
    <h1>Page 2</h1>
    <ul>
        <li><a href="index.html">Index Page</a></li>
        <li><a href="page1.html">Page 1</a></li>
        <li><a href="https://example.com">Example Site</a></li>
        <li><a href="https://www.google.com">Google</a></li>
    </ul>
</body>
</html>`)

var mock = []byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Empty HTML Document</title>
</head>
<body>
    
</body>
</html>`)

func TestWgetParser(t *testing.T) {

	wgReqs := []Wget{
		{
			LinksVisited: make(map[string]struct{}, 10),
			URL:          "https://test.com/index.html",
			OnlyInternal: true,
			Links:        []string{"https://test.com/index.html"},
		},
		{
			LinksVisited: make(map[string]struct{}, 10),
			URL:          "https://test.com/index.html",
			Depth:        1,
			Links:        []string{"https://test.com/index.html"},
		},
	}

	testAnser := [][]string{
		{"https://test.com/index.html", "https://test.com/page2.html", "https://test.com/page1.html"},
		{"https://test.com/index.html", "https://www.google.com", "https://example.com", "https://test.com/page2.html", "https://test.com/page1.html"},
	}

	mp := map[string][]byte{
		"https://test.com/index.html": page1,
		"https://test.com/page1.html": page2,
		"https://test.com/page2.html": page3,
		"https://example.com":         mock,
		"https://www.google.com":      mock,
	}

	for i, v := range wgReqs {
		ans, err := Mock(&v, mp)
		if err != nil {
			t.Errorf(err.Error())
		}
		assert.Equal(t, testAnser[i], ans, "not correct")
	}

}

func Mock(wg *Wget, mp map[string][]byte) ([]string, error) {
	allStrings := make([]string, 0, 100)
	allStrings = append(allStrings, wg.URL)
	var needIncr bool
	if wg.Depth == -1 {
		wg.Depth = 0
		needIncr = true
	}
	for i := 0; i <= wg.Depth; i++ {
		if len(wg.Links) == 0 {
			return nil, nil
		}
		nextDepthLinks := make([]string, 0, len(wg.Links))

		for _, v := range wg.Links {

			wg.LinksVisited[v] = struct{}{}
			htmlData, ok := mp[v]
			if !ok {
				fmt.Printf("key %s ", v)
			}

			links, err := wg.ParseHTML(v, htmlData)
			if err != nil {
				return nil, err
			}

			nextDepthLinks = append(nextDepthLinks, links...)
			allStrings = append(allStrings, links...)
			wg.FilesDownloaded++
		}
		wg.Links = nextDepthLinks
		if needIncr {
			wg.Depth++
		}
	}
	return allStrings, nil
}
