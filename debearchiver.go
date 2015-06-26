package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://eksisozluk.com/debe", nil)
	if err != nil {
		fmt.Println("error", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:36.0) Gecko/20100101 Firefox/36.0")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error", err)
	}
	defer res.Body.Close()

	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error", err)
	}

	parseHtml(strings.NewReader(string(contents)))
}

func parseHtml(r io.Reader) {
	parseEnable := false
	tokenT := ""

	d := html.NewTokenizer(r)
	for {
		// token type
		tokenType := d.Next()
		if tokenType == html.ErrorToken {
			return
		}
		token := d.Token()
		switch tokenType {
		case html.StartTagToken: // <tag>
			// type Token struct {
			//     Type     TokenType
			//     DataAtom atom.Atom
			//     Data     string
			//     Attr     []Attribute
			// }
			//
			// type Attribute struct {
			//     Namespace, Key, Val string
			// }
			if token.Data == "ol" {
				parseEnable = true
			}

			if parseEnable {
				if token.Data == "a" && len(token.Attr) > 0 {
					fmt.Println(token.Attr[0].Val)
				}
			}
			tokenT = token.Data

		case html.TextToken: // text between start and end tag
			if parseEnable && (tokenT == "span" || tokenT == "div") {
				fmt.Println(token.Data)
			}
		case html.EndTagToken: // </tag>
			tokenT = ""
			if token.Data == "ol" {
				parseEnable = false
			}
		case html.SelfClosingTagToken: // <tag/>
		}
	}
}
