package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Entry struct {
	url     string
	subject string
	author  string
}

var aylar = [...]string{
	"Ocak",
	"Şubat",
	"Mart",
	"Nisan",
	"Mayıs",
	"Haziran",
	"Temmuz",
	"Ağustos",
	"Eylül",
	"Ekim",
	"Kasım",
	"Aralık",
}

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

	var entries []Entry
	parseHtml(strings.NewReader(string(contents)), &entries)

	writeFile(entries)
}

func writeFile(entries []Entry) {
	year, month, day := time.Now().Date()
	filename := fmt.Sprintf("%d-%02d-%02d-%d-%s-%d-eksisozluk-debe.md", year, int(month), day, day, aylar[int(month)-1], year)

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString(fmt.Sprintf("---\nlayout: post\ntitle: %d %s %d Ekşi Sözlük Debe\n---\n\n", day, aylar[int(month)-1], year))
	if err != nil {
		panic(err)
	}

	for _, v := range entries {
		_, err = f.WriteString(fmt.Sprintf("* [%s](http://eksisozluk.com/%s)\n", v.subject, v.url))
		if err != nil {
			panic(err)
		}
	}
}

func parseHtml(r io.Reader, entries *[]Entry) {
	parseEnable := false
	tokenT := ""
	var entry Entry

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
					entry.author = ""
					entry.subject = ""
					entry.url = token.Attr[0].Val
				}
			}
			tokenT = token.Data

		case html.TextToken: // text between start and end tag
			if parseEnable && (tokenT == "span" || tokenT == "div") {
				if len(entry.subject) > 0 {
					entry.author = token.Data
					*entries = append(*entries, entry)
				} else {
					entry.subject = token.Data
				}
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
