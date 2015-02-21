package api
import "testing"


func TestGetTcode(t *testing.T) {
	tcode,_ := GetTcode("n9902bn")
	if tcode != "399863" {
		t.Errorf("get tcode error %s", tcode)
	}
}

