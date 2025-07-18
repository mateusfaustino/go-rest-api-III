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
	_ "github.com/mateusfaustino/go-rest-api-III/internal/entity"
	"github.com/mateusfaustino/go-rest-api-III/internal/infra/database"
	entityPkg "github.com/mateusfaustino/go-rest-api-III/pkg/entity"
	"gorm.io/gorm"
)

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserDb       database.UserInterface
	RoleDB       database.RoleInterface
	Jwt          *jwtauth.JWTAuth
	JwtExpiresIn int
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func NewUserHandler(db database.UserInterface, roleDB database.RoleInterface, jwt *jwtauth.JWTAuth, jwtExpiresIn int) *UserHandler {
	return &UserHandler{
		UserDb:       db,
		RoleDB:       roleDB,
		Jwt:          jwt,
		JwtExpiresIn: jwtExpiresIn,
	}
}

// Instância do validador globalmente
var validate = validator.New()

// GetJWT godoc
// @Summary: Get a JWT token
// @Description: Get a JWT token with the given email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.GetJWTInput true "User data"
// @Success 200 {object} dto.GetJWTOutput
// @Failure 400 {object} Error
// @Failure 401 {object} Error
// @Failure 500 {object} Error
// @Router /auth/login [post]
func (uh *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "request body is empty"})
		return
	}

	defer r.Body.Close()

	var userInput dto.GetJWTInput
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON format"})
		return
	}

	if strings.TrimSpace(userInput.Email) == "" || strings.TrimSpace(userInput.Password) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "email and password are required"})
		return
	}

	u, err := uh.UserDb.FindUserByEmail(userInput.Email)
	if err != nil {
		time.Sleep(500 * time.Millisecond) // Delay to prevent timing attacks
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
		return
	}

	if !u.ValidatePassword(userInput.Password) {
		time.Sleep(500 * time.Millisecond) // Delay to prevent timing attacks
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
		return
	}

	// Verifica se o JWT está configurado corretamente
	if uh.Jwt == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "JWT service not properly configured"})
		return
	}

	// Verifica se o tempo de expiração é válido
	if uh.JwtExpiresIn <= 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JWT expiration time"})
		return
	}

	// Cria o payload do token
	claims := map[string]interface{}{
		"sub":  u.ID.String(),
		"exp":  time.Now().Add(time.Second * time.Duration(uh.JwtExpiresIn)).Unix(),
		"role": u.RoleID,
	}

	// Gera o token com timeout
	tokenChan := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		_, tokenString, err := uh.Jwt.Encode(claims)
		if err != nil {
			errChan <- err
			return
		}
		tokenChan <- tokenString
	}()

	// Aguarda a geração do token com timeout
	select {
	case tokenString := <-tokenChan:
		accessToken := dto.GetJWTOutput{AccessToken: tokenString}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(accessToken)
	case err := <-errChan:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("could not generate token: %v", err)})
	case <-time.After(5 * time.Second):
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "token generation timeout"})
	}
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

// Create User
// @Summary: Create a new user
// @Description: Create a new user with the given name, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.CreateUserInput true "User data"
// @Success 201 {object} dto.CreateUserInput
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /auth/register [post]
func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userInput dto.CreateUserInput

	// Fechar o corpo da requisição após uso
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&userInput)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": %s}`, err.Error()), http.StatusBadRequest)
		return
	}

	roleDB := uh.RoleDB

	roleCustomer, err := roleDB.FindRoleByName("customer")

	if err != nil {
		http.Error(w, `{"error": "could not find customer role"}`, http.StatusInternalServerError)
		return
	}

	u, err := entity.NewUser(userInput.Name, userInput.Email, userInput.Password, roleCustomer.ID)

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

// GetUserById godoc
// @Summary: Get a user by ID
// @Description: Retrieve user information by ID
// @Tags user
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} entity.User
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /user/{id} [get]
// @Security ApiKeyAuth
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

// UpdateOwnProfile godoc
// @Summary: Update authenticated user profile
// @Description: Update the profile of the logged in user
// @Tags user
// @Accept json
// @Produce json
// @Param profile body dto.UpdateOwnProfileInput true "Profile data"
// @Success 200 {object} dto.UpdateOwnProfileInput
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 500 {object} Error
// @Router /user/profile [put]
// @Security ApiKeyAuth
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
	// foundedUser.Role = userRole

	err = uh.UserDb.UpdateUser(foundedUser)

	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": %s}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userInput)

}

// ShowOwnProfile godoc
// @Summary: Show authenticated user profile
// @Description: Retrieve the profile of the logged in user
// @Tags user
// @Produce json
// @Success 200 {object} entity.User
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /user/profile [get]
// @Security ApiKeyAuth
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
