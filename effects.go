package main

import (
	"encoding/json"
	"net/http"
)

type Effect struct {
	Key  string `json:"key"`
	Text string `json:"text"`
}

func randomMagicalEffect(w http.ResponseWriter, r *http.Request) {
	effect, err := randomEffect()
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	json.NewEncoder(w).Encode(effect)
}
