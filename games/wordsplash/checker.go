package wordsplash

import (
	"encoding/json"
	"net/http"
)

type result struct {
	Found int `json:"found"`
}

func wordExists(word string) bool {
	resp, _ := http.Get("http://www.anagramica.com/lookup/")
	var v result
	json.NewDecoder(resp.Body).Decode(v)
	return v.Found > 0
}
