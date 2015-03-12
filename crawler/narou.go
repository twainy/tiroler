package crawler

import (
	"github.com/PuerkitoBio/goquery"
    "fmt"
    "log"
    "regexp"
)

type Novel struct {
    tcode string
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
            n.tcode = tcode

        } else {
        }
    })
    return n, err
}
