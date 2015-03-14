package crawler
import (
    "testing"
)


func TestGetnovel(t *testing.T) {
    n,_ := GetNovel("n9669bk")
    if n.Tcode != "369633" {
        t.Error("invalid tcode")
    }
    if n.ContentList[0].Ctype != Chapter {
        t.Error("invalid content")
    }
    if n.ContentList[0].Text != "第１章　幼年期" {
        t.Error("invalid text")
    }
    if n.ContentList[1].Ctype != Sublist {
        t.Error("invalid content")
    }
    if n.ContentList[1].Text != "プロローグ" {
        t.Error("invalid text")
    }

}
