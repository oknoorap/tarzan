package request

import (
"reflect"
"sync"
"net/url"
"encoding/xml"
"github.com/speps/go-hashids"
"strings"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Items RSSItems `xml:"channel"`
}

type RSSItems struct {
	XMLName xml.Name `xml:"channel"`
	ItemList []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title string `xml:"title"`
	Url string `xml:"link"`
	Desc string `xml:"description"`
}

type Bing struct {
	Items []BingItems `json:"items"`
}

func (bing *Bing) HasItem() bool {
	if bing.Items != nil {
		return true
	}

	return false
}

func (bing *Bing) Size() int {
	return len(bing.Items)
}

func (bing *Bing) Add(item BingItems) {
	bing.Items = append(bing.Items, item)
}

type BingItems struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Desc string `json:"desc"`
	Author string `json:"author"`
	Size int64 `json:"size"`
	Url string `json:"url"`
	Web string `json:"web"`
}

func getAuthor(uri string) (string, string) {
	u, _ := url.Parse(uri)
	originhost := u.Host
	// Replace all web 2.0
	switch true {
		case strings.Contains(originhost, ".files.wordpress"):
			originhost = strings.Replace(originhost, ".files.", ".", -1)
	}
	
	host := strings.Replace(originhost, "-", " ", -1)
	host = strings.Replace(host, "www.", "", -1)
	host = strings.Replace(host, "www2.", "", -1)
	host = strings.Replace(host, "cdn.", "", -1)

	author := func(h string, tld string) string {
		h = strings.Replace(h, "." + tld, "", -1)
		shost := strings.Split(h, ".")
		h = strings.Join(shost, " ")
		h = strings.Title(h)
		return h
	}

	switch true {
		case strings.HasSuffix(host, ".com"):
			host = author(host, "com")
		case strings.HasSuffix(host, ".net"):
			host = author(host, "net") + " Network"
		case strings.HasSuffix(host, ".edu"):
			host = "Institute of " + author(host, "edu")
		case strings.HasSuffix(host, ".org"):
			host = author(host, "org") + " Organization"
		case strings.HasSuffix(host, ".gov"):
			host = author(host, "gov") + " U.S. government agency"
		case strings.HasSuffix(host, ".mil"):
			host = "U.S. Military:" + author(host, "mil")
	}

	return host, u.Scheme +"://"+ originhost
}

func (bing *Bing) GetResult() []BingItems {
	return bing.Items
}

func (bing *Bing) Search(query, count string) {
	query = Humanize(query, false)
	query = strings.ToLower(query)
	query = url.QueryEscape(query)
	tor := &Tor{}
	query = "http://www.bing.com/search?q="+ query +"&count="+ count +"&format=rss"
	body, _ := tor.Open(query, "http://www.bing.com/search?q="+ query)
	//body := ``

	if body != "" {
		var rss RSS
		err := xml.Unmarshal([]byte(body), &rss)
		bing_result := rss.Items.ItemList

		if err == nil && bing_result != nil {
			// Generate hashId
			hd := hashids.NewData()
			hd.MinLength = 5

			// Parse results interface
			res := reflect.ValueOf(bing_result)

			// Get File Size
			size := make(chan int64)
			defer close(size)
			var done sync.WaitGroup
			for i := 0; i < res.Len(); i++ {
				go func(uri string) {
					size <- GetUrlSize(uri)
				}(res.Index(i).FieldByName("Url").String())
			}
			done.Wait()

			for i := 0; i < res.Len(); i++ {
				me := res.Index(i)

				// Generate hash id by url
				hd.Salt = me.FieldByName("Url").String()
				h := hashids.NewWithData(hd)
				id, _ := h.Encode([]int{17,1})

				// Check title
				final_title := TitleCase(me.FieldByName("Title").String())
				if final_title == "" {
					final_title = "Untitled document"
				}

				// Check description
				final_desc := CleanDesc(me.FieldByName("Desc").String())
				if final_desc == "" {
					final_desc = "No description available."
				}

				author, web := getAuthor(me.FieldByName("Url").String())
				
				// Append results
				bing.Add(BingItems{
					Id: id,
					Title: final_title,
					Desc: final_desc,
					Url: me.FieldByName("Url").String(),
					Author: author,
					Web: web,
					Size: <-size,
				})
			}
		}
	}
}