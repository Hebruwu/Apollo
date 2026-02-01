package api

import (
	"net/http"

	"apollo.io/api/users"
	"apollo.io/utils"
)

func NewAPIV1(config utils.ServerConfig) http.Handler {
	api_v1 := http.NewServeMux()

	api_v1.Handle("/users", http.StripPrefix("/users", users.NewUsersBase(config)))

	return api_v1
}
