package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-playground/validator"
	"github.com/mateusfaustino/go-rest-api-III/internal/dto"
	"github.com/mateusfaustino/go-rest-api-III/internal/entity"
	"github.com/mateusfaustino/go-rest-api-III/internal/infra/database"
	entityPkg "github.com/mateusfaustino/go-rest-api-III/pkg/entity"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserDb       database.UserInterface
	Jwt          *jwtauth.JWTAuth
	JwtExpiresIn int
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func NewUserHandler(db database.UserInterface, jwt *jwtauth.JWTAuth, jwtExpiresIn int) *UserHandler {
	return &UserHandler{
		UserDb:       db,
		Jwt:          jwt,
		JwtExpiresIn: jwtExpiresIn,
	}
}

// Instância do validador globalmente
var validate = validator.New()

func (uh *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var userInput dto.GetJWTInput

	if r.Body == nil {
		http.Error(w, `{"error": "request body is empty"}`, http.StatusBadRequest)
		return
	}

	// Fechar o corpo da requisição após uso
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&userInput)

	if strings.TrimSpace(userInput.Email) == "" || strings.TrimSpace(userInput.Password) == "" {
		http.Error(w, `{"error": "email and password are required"}`, http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, `{"error": "invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	u, err := uh.UserDb.FindUserByEmail(userInput.Email)

	if err != nil || !u.ValidatePassword(userInput.Password) {
		time.Sleep(500 * time.Millisecond) // Pequeno delay para evitar timing attacks
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	_, tokenString, err := uh.Jwt.Encode(map[string]interface{}{
		"sub":  u.ID.String(),
		"exp":  time.Now().Add(time.Second * time.Duration(uh.JwtExpiresIn)).Unix(),
		"role": u.Role,
	})

	if err != nil {
		http.Error(w, `{"error": "could not generate token"}`, http.StatusInternalServerError)
		return
	}

	accessToken := AccessTokenResponse{AccessToken: tokenString}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)

}

func (uh *UserHandler) TestManager(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(`{message: ok}`)
}

func (uh *UserHandler) TestCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(`{message: ok}`)
}

func (uh *UserHandler) TestAdmin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(`{message: ok}`)
}

func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userInput dto.CreateUserInput

	// Fechar o corpo da requisição após uso
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&userInput)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": %s}`, err.Error()), http.StatusBadRequest)
		return
	}

	u, err := entity.NewUser(userInput.Name, userInput.Email, userInput.Password, "customer")

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": %s}`, err.Error()), http.StatusBadRequest)
		return
	}

	err = uh.UserDb.CreateUser(u)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": %s}`, err.Error()), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "user created successfully"})

}

func (uh *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, `{"error": "missing user ID"}`, http.StatusBadRequest)
		return
	}

	user, err := uh.UserDb.FindUserById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}

func (uh *UserHandler) UpdateOwnProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close() // Garante que qualquer recurso seja fechado

	var userInput dto.UpdateOwnProfileInput
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// Valida os campos obrigatórios
	if err := validate.Struct(userInput); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorsMap := make(map[string]string)

		for _, fieldErr := range validationErrors {
			errorsMap[fieldErr.Field()] = fmt.Sprintf("Field '%s' is required and must be valid", fieldErr.Field())
		}

		response, _ := json.Marshal(map[string]interface{}{"errors": errorsMap})
		http.Error(w, string(response), http.StatusBadRequest)
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())

	userId, ok := claims["sub"].(string)
	if !ok || userId == "" {
		http.Error(w, `{"error": "invalid token: missing 'sub'"}`, http.StatusForbidden)
		return
	}

	userRole, ok := claims["role"].(string)
	if !ok || userRole == "" {
		userRole = "customer"
		// http.Error(w, `{"error": "invalid token: missing 'role'"}`, http.StatusForbidden)
		// return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": %s}`, err.Error()), http.StatusInternalServerError)
		return
	}

	if userInput.NewPassword == "" {
		userInput.NewPassword = userInput.Password
	}

	hashedPassword, err := entityPkg.HashPassword(userInput.NewPassword)

	if err != nil {
		http.Error(w, `{"error": "failed to hash password"}`, http.StatusInternalServerError)
		return
	}

	foundedUser, err := uh.UserDb.FindUserById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	if err != nil || !foundedUser.ValidatePassword(userInput.Password) {
		time.Sleep(500 * time.Millisecond) // Pequeno delay para evitar timing attacks
		http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	foundedUserByEmail, _ := uh.UserDb.FindUserByEmail(userInput.Email)
	if foundedUserByEmail != nil {
		if foundedUserByEmail.ID != foundedUser.ID {
			http.Error(w, `{"error": "this email is already used"}`, http.StatusBadRequest)
			return
		}
	}

	foundedUser.Name = userInput.Name
	foundedUser.Email = userInput.Email
	foundedUser.Password = hashedPassword
	foundedUser.Role = userRole

	err = uh.UserDb.UpdateUser(foundedUser)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": %s}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userInput)

}

func (uh *UserHandler) ShowOwnProfile(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	userId, ok := claims["sub"].(string)

	if !ok {
		http.Error(w, `{"error": "invalid token"}`, http.StatusForbidden)
		return
	}

	userFound, err := uh.UserDb.FindUserById(userId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userFound)

}
