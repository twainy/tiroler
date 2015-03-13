package crawler
import (
    "testing"
)


func TestGetnovel(t *testing.T) {
    n,_ := GetNovel("n9669bk")
    if n.tcode != "369633" {
        t.Error("invalid tcode")
    }
    if n.content_list[0].ctype != Chapter {
        t.Error("invalid content")
    }
    if n.content_list[0].text != "第１章　幼年期" {
        t.Error("invalid text")
    }
    if n.content_list[1].ctype != Sublist {
        t.Error("invalid content")
    }
    if n.content_list[1].text != "プロローグ" {
        t.Error("invalid text")
    }

}
