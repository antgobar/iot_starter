package api

import (
	"errors"
	"iotstarter/internal/config"
	"iotstarter/internal/model"
	"log"
	"net/http"
)

func getUserFromRequest(r *http.Request) (*model.User, error) {
	val := r.Context().Value(config.UserKey)
	user, ok := val.(*model.User)
	if !ok || user == nil {
		return nil, errors.New("no user in context")
	}
	log.Println("User in request", user)
	return user, nil
}
