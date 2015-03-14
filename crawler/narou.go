package crawler

import (
	"github.com/PuerkitoBio/goquery"
    "fmt"
    "log"
    "regexp"
    "strconv"
)

type NovelContentType int

type Novel struct {
    Tcode string
    ContentList []NovelContent
}

const ( // NovelContent Type
    Chapter NovelContentType = iota
    Sublist
)
type NovelContent struct {
    Ctype NovelContentType
    Text string
    SublistId int
}

func GetNovel(ncode string) (Novel, error) {
    doc, err := goquery.NewDocument(fmt.Sprintf("http://ncode.syosetu.com/%s/", ncode));
    if err != nil {
       log.Fatal(err);
        
    }
    
    n := Novel{}
    
    doc.Find("#novel_footer ul li").Each( func(i int, s *goquery.Selection) {
        if s.Find("a").Text() == "TXTダウンロード" {
            href,_ := s.Find("a").Attr("href")
            re, _ := regexp.Compile("[0-9]{6}")
            tcode := string(re.Find([]byte(href)))
            n.Tcode = tcode

        } else {
        }
    })
    doc.Find("div.index_box").Children().Each(func(i int, s *goquery.Selection) {
        if s.HasClass("chapter_title") {
            c := NovelContent{}
            c.Ctype = Chapter
            c.Text = s.Text()
            n.ContentList = append(n.ContentList, c)
        }
        if s.HasClass("novel_sublist2") {
            subtitle := s.Find(".novel_sublist2 dd.subtitle")
            url,_ := s.Find(".novel_sublist2 a").Attr("href")
            re, _ := regexp.Compile("/[0-9]+/")
            sublist_id,_ := strconv.Atoi(string(re.Find([]byte(url))))
            c := NovelContent{}
            c.Ctype = Sublist
            c.Text = subtitle.Text()
            c.SublistId = sublist_id
            n.ContentList = append(n.ContentList, c)
        }
    });
    return n, err
}
