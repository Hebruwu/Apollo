package apolloapiv1

import (
	"net/http"

	"apollo.io/apolloapiv1/users"
	"apollo.io/serverutils"
)

func NewAPIV1(config serverutils.ServerConfig) http.Handler {
	apiV1 := http.NewServeMux()

	apiV1.Handle("/users/", http.StripPrefix("/users", users.NewUsersBase(config)))

	return apiV1
}
