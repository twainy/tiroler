package api

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
);
const API_BASE_URL = "https://www.kimonolabs.com/api/ck52w24s"
const API_KEY_TCODE = "tEqeo1dz9ZfMuUvCeR55gM80kT6AkJzX"

type TcodeResponse struct {
	Tcode string `json:"tcode"`
	Status string `json:"status"`
}

func GetTcode(ncode string) (string, error) {
	url := API_BASE_URL + "?apikey="+API_KEY_TCODE+"&ncode="+ncode+"&kimmodify=1"
	resp, err := http.Get(url)
	if(err == nil){
		body, err := ioutil.ReadAll(resp.Body)
		if(err == nil) {
			var data TcodeResponse
			json.Unmarshal(body, &data)

			return data.Tcode,err
		} else {
		}
	} else {
	}
	return "",err
}
// get request to API

