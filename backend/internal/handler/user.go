package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"EMS/internal/cognito"
	"EMS/internal/middleware"
	"EMS/internal/model"
	"EMS/internal/service"
	"EMS/internal/utils"
)

type UserHandler struct {
	service  *service.UserService
	verifier *cognito.JWTVerifier
}

func NewUserHandler(service *service.UserService, verifier *cognito.JWTVerifier) *UserHandler {
	return &UserHandler{
		service:  service,
		verifier: verifier,
	}
}

// GET /api/me
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	sub, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok || sub == "" {
		utils.SendJSON(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Unauthorized",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	user, err := h.service.GetUserBySub(ctx, sub)
	if err != nil && !errors.Is(err, service.ErrUserNotFound) {
		utils.SendJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to query user: " + err.Error(),
		})
		return
	}

	if user != nil {
		utils.SendJSON(w, http.StatusOK, utils.APIResponse{
			Status:  "success",
			Message: "User fetched successfully",
			Data:    user,
		})
		return
	}

	// User not found, try to create from ID token
	idTokenStr := r.Header.Get("X-Id-Token")
	if idTokenStr == "" {
		utils.SendJSON(w, http.StatusBadRequest, utils.APIResponse{
			Status:  "error",
			Message: "User not found and X-Id-Token is missing",
		})
		return
	}

	idToken, err := h.verifier.VerifyIDToken(idTokenStr)
	if err != nil {
		utils.SendJSON(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Invalid ID Token: " + err.Error(),
		})
		return
	}

	claims, ok := idToken.Claims.(jwt.MapClaims)
	if !ok {
		utils.SendJSON(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Invalid ID Token claims",
		})
		return
	}

	// Validate sub matches
	idTokenSub, _ := claims["sub"].(string)
	if idTokenSub != sub {
		utils.SendJSON(w, http.StatusUnauthorized, utils.APIResponse{
			Status:  "error",
			Message: "Sub mismatch between Access and ID tokens",
		})
		return
	}

	email, _ := claims["email"].(string)
	name, _ := claims["name"].(string)
	if name == "" {
		name, _ = claims["cognito:username"].(string)
	}

	newUser := model.User{
		Sub:   sub,
		Email: email,
		Name:  name,
	}

	createdUser, err := h.service.CreateUser(ctx, newUser)
	if err != nil {
		if errors.Is(err, service.ErrDuplicateEmail) || errors.Is(err, service.ErrDuplicateSub) {
			utils.SendJSON(w, http.StatusBadRequest, utils.APIResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}
		utils.SendJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to create user: " + err.Error(),
		})
		return
	}

	utils.SendJSON(w, http.StatusCreated, utils.APIResponse{
		Status:  "success",
		Message: "User provisioned successfully",
		Data:    createdUser,
	})
}

// GET /api/users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	role, _ := r.Context().Value(middleware.UserRoleContextKey).(string)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	users, err := h.service.GetAllUsers(ctx, role)
	if err != nil {
		if errors.Is(err, service.ErrAdminRequired) {
			utils.SendJSON(w, http.StatusForbidden, utils.APIResponse{
				Status:  "error",
				Message: "Admin access required",
				Data:    nil,
			})
			return
		}
		utils.SendJSON(w, http.StatusInternalServerError, utils.APIResponse{
			Status:  "error",
			Message: "Failed to fetch users: " + err.Error(),
			Data:    nil,
		})
		return
	}

	utils.SendJSON(w, http.StatusOK, utils.APIResponse{
		Status:  "success",
		Message: "Users fetched successfully",
		Data:    users,
	})
}
