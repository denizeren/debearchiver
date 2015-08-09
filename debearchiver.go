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
	id      string
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
	var entries []Entry
	parseDebeHtml(strings.NewReader(makeHttpRequest("https://eksisozluk.com/debe")), &entries)
	time.Sleep(1 * time.Second)

	writeFile(entries)
}

func makeHttpRequest(url string) string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
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

	return string(contents)
}

func writeFile(entries []Entry) {
	year, month, day := time.Now().Date()
	filename := fmt.Sprintf("%d-%02d-%02d-%d-%s-%d-eksisozluk-debe.md", year, int(month), day, day, aylar[int(month)-1], year)

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString(fmt.Sprintf("---\nlayout: post\ntitle: %d %s %d Ekşi Sözlük Debe\ndata:\n", day, aylar[int(month)-1], year))
	if err != nil {
		panic(err)
	}

	for _, v := range entries {
		_, err = f.WriteString(fmt.Sprintf("- entry_name: |\n    %s\n  entry_id: %s\n", v.subject, v.id))
		if err != nil {
			panic(err)
		}

		entryData := parseEntryHtml(strings.NewReader(makeHttpRequest(fmt.Sprintf("https://eksisozluk.com/entry/%s", v.id))))
		_, err = f.WriteString(fmt.Sprintf("  entry_content: |\n    %s\n  entry_writer: %s\n", entryData, v.author))
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}

	_, err = f.WriteString("---\n")
	if err != nil {
		panic(err)
	}
}

func parseDebeHtml(r io.Reader, entries *[]Entry) {
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
					entryList := strings.Split(entry.url, "%23")
					entry.id = entryList[len(entryList)-1]
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

func parseEntryHtml(r io.Reader) string {
        parseEnable := false
	retVal := ""

        d := html.NewTokenizer(r)
        for {
                // token type
                tokenType := d.Next()
                if tokenType == html.ErrorToken {
			return retVal
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
                        if parseEnable == true {
                                retVal += string(d.Raw())
                        }
                        if token.Data == "div" {
                                if len(token.Attr) > 0 && token.Attr[0].Key == "class" && token.Attr[0].Val == "content" {
                                        parseEnable = true
                                }
                        }

                case html.TextToken: // text between start and end tag
                        if parseEnable == true {
                                retVal += string(d.Raw())
                        }

                case html.EndTagToken: // </tag>
                        if token.Data == "div" {
                                parseEnable = false
                        }
                        if parseEnable == true {
                                retVal += string(d.Raw())
                        }
                case html.SelfClosingTagToken: // <tag/>
                        if parseEnable == true {
                                retVal += string(d.Raw())
                        }
                }
        }

	return retVal
}
