package main

import (
	"testing"
	"encoding/json"
)

func TestMain(t *testing.T) {
	mapA := map[string]string{"aaa":"bbbb","ccc":"ddd"}
	mapB, _ := json.Marshal(mapA)
	if string(mapB) != "{\"aaa\":\"bbbb\",\"ccc\":\"ddd\"}" {
        t.Error("json marshal error")
        
    }
    
}
