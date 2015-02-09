package tiroler

import (
	"testing"
	"encoding/json"
	"fmt"
)

func TestMain(t *testing.T) {
	mapA := map[string]string{"aaa":"bbbb","ccc":"ddd"}
	mapB, _ := json.Marshal(mapA)
	fmt.Println(string(mapB))
}
