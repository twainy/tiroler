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
    Tcode string `json:"tcode"`
    ContentList []NovelContent `json:"content_list"`
}

const ( // NovelContent Type
    Chapter NovelContentType = iota
    Sublist
)
type NovelContent struct {
    Ctype NovelContentType `json:"ctype"`
    Text string `json:"text"`
    SublistId int `json:"sublist_id"`
    Content string `json:"content"`
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
            re, _ := regexp.Compile("/([0-9]+)/")
            sublist_id,_ := strconv.Atoi(re.FindStringSubmatch(url)[1])
            c := NovelContent{}
            c.Ctype = Sublist
            c.Text = subtitle.Text()
            c.SublistId = sublist_id
            n.ContentList = append(n.ContentList, c)
        }
    });
    return n, err
}

func GetNovelContent(ncode string ,chapter_id int) NovelContent {
    url := fmt.Sprintf("http://ncode.syosetu.com/%s/%d", ncode, chapter_id)
    doc, err := goquery.NewDocument(url)
    log.Println("get content " + url)
    if err != nil {
       log.Fatal(err);
    }
    content_title := doc.Find(".novel_subtitle").Text()
    content := doc.Find("#novel_honbun").Text()

    c := NovelContent{}
    c.Text = content_title
    c.Content = content
    return c
}
