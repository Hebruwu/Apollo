package users

import (
	"encoding/json"
	"net/http"

	"apollo.io/utils"
)

type userRoute struct {
	serverConf *utils.ServerConfig
}

func NewUsersBase(serverConfig utils.ServerConfig) http.Handler {
	users := http.NewServeMux()
	ur := userRoute{serverConf: &serverConfig}

	users.HandleFunc("POST /login", ur.authenticate)
	users.HandleFunc("POST /register", ur.createUser)
	users.HandleFunc("POST /logout", ur.logoutUser)
	users.HandleFunc("POST /refresh", ur.refreshToken)
	users.HandleFunc("PUT /{id}", ur.updateUser)
	users.HandleFunc("DELETE /{id}", ur.deleteUser)

	return users
}

func (ur userRoute) authenticate(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement getUser
}

// Request body for user creation.
type NewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// Handles user creation.
func (ur userRoute) createUser(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	var user NewUser
	err := json.NewDecoder(r.Body).Decode(&user)

	// Malformed payload.
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Look into https://snyk.io/blog/secure-password-hashing-in-go/ to generate secure password hash.

}

func (ur userRoute) updateUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement updateUser
}

func (ur userRoute) deleteUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement deleteUser
}

func (ur userRoute) logoutUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement logoutUser
}

func (ur userRoute) refreshToken(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement refreshToken
}
