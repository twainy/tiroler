package crawler
import "testing"


func TestGetnovel(t *testing.T) {
    n,_ := GetNovel("n9902bn")
    if n.tcode != "399863" {
        t.Error("invalid tcode")
    }
}
