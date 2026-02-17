package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"apollo.io/objects/request"
	"apollo.io/objects/response"
	"apollo.io/objects/shared"
	"apollo.io/services"
	"apollo.io/utils"
)

type userRoute struct {
	services.UserService
	utils.ServerConfig
}

func NewUsersBase(serverConfig utils.ServerConfig) http.Handler {
	users := http.NewServeMux()
	ur := userRoute{
		services.NewUserService(serverConfig.PostgresConnection, serverConfig.Logger),
		serverConfig,
	}

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

// Handles user creation.
func (ur userRoute) createUser(w http.ResponseWriter, r *http.Request) {
	var user request.NewUser
	err := json.NewDecoder(r.Body).Decode(&user)

	// Malformed payload.
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = ur.CreateUser(r.Context(), user.Username, user.Email, user.Password)
	
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, shared.ErrUsernameAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(response.StatusResponse{Error: response.USERNAME_ALREADY_EXISTS})
		} else {
			ur.Logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.StatusResponse{Error: response.UNEXPECTED_ERROR})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.StatusResponse{Success: response.USER_CREATED})
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
