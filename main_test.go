package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	mapA := map[string]string{"aaa": "bbbb", "ccc": "ddd"}
	mapB, _ := json.Marshal(mapA)
}
