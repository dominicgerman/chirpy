package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

var ErrNotExist = errors.New("resource does not exist")

func (cfg apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "No auth header provided")
		return
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		respondWithError(w, http.StatusUnauthorized, "Couldn't parse auth header")
		return
	}

	if splitAuth[1] != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Incorrect API key")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, "")
		return
	}

	_, err = cfg.DB.UpgradeUser(params.Data.UserID)
	if err != nil {
		if !errors.Is(err, ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't upgrade user")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")

}
