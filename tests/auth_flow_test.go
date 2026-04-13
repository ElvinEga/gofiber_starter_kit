package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/ElvinEga/adeya_backend/config"
	"github.com/ElvinEga/adeya_backend/database"
	"github.com/ElvinEga/adeya_backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type authPayload struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"user"`
}

type errorPayload struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func setupAuthTestApp(t *testing.T) *fiber.App {
	t.Helper()

	t.Setenv("DATABASE_URL", "")
	t.Setenv("DB_PATH", "file::memory:?cache=shared")
	t.Setenv("JWT_SECRET", "test-secret")

	config.InitConfig()
	database.ConnectDB()
	database.MigrateDB()

	app := fiber.New()
	routes.SetupRoutes(app)
	return app
}

func performJSONRequest(t *testing.T, app *fiber.App, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	recorder.Code = resp.StatusCode
	_, _ = recorder.Body.ReadFrom(resp.Body)
	return recorder
}

func TestRegisterReturnsTokenPair(t *testing.T) {
	app := setupAuthTestApp(t)

	resp := performJSONRequest(t, app, "POST", "/api/auth/register", map[string]string{
		"name":     "Test User",
		"email":    "register@example.com",
		"password": "Password123!",
	})

	require.Equal(t, 201, resp.Code)

	var payload authPayload
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	assert.Equal(t, "success", payload.Status)
	assert.NotEmpty(t, payload.AccessToken)
	assert.NotEmpty(t, payload.RefreshToken)
	assert.Equal(t, "register@example.com", payload.User.Email)
}

func TestRegisterRejectsDuplicateEmail(t *testing.T) {
	app := setupAuthTestApp(t)

	performJSONRequest(t, app, "POST", "/api/auth/register", map[string]string{
		"name":     "Test User",
		"email":    "duplicate@example.com",
		"password": "Password123!",
	})

	resp := performJSONRequest(t, app, "POST", "/api/auth/register", map[string]string{
		"name":     "Another User",
		"email":    "duplicate@example.com",
		"password": "Password123!",
	})

	require.Equal(t, 409, resp.Code)

	var payload errorPayload
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	assert.Equal(t, "error", payload.Status)
}

func TestLoginReturnsTokenPair(t *testing.T) {
	app := setupAuthTestApp(t)

	performJSONRequest(t, app, "POST", "/api/auth/register", map[string]string{
		"name":     "Login User",
		"email":    "login@example.com",
		"password": "Password123!",
	})

	resp := performJSONRequest(t, app, "POST", "/api/auth/login", map[string]string{
		"email":    "login@example.com",
		"password": "Password123!",
	})

	require.Equal(t, 200, resp.Code)

	var payload authPayload
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	assert.Equal(t, "success", payload.Status)
	assert.NotEmpty(t, payload.AccessToken)
	assert.NotEmpty(t, payload.RefreshToken)
}

func TestLoginRejectsInvalidCredentials(t *testing.T) {
	app := setupAuthTestApp(t)

	resp := performJSONRequest(t, app, "POST", "/api/auth/login", map[string]string{
		"email":    "missing@example.com",
		"password": "Password123!",
	})

	require.Equal(t, 401, resp.Code)

	var payload errorPayload
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	assert.Equal(t, "error", payload.Status)
}

func TestRefreshReturnsNewTokenPair(t *testing.T) {
	app := setupAuthTestApp(t)

	registerResp := performJSONRequest(t, app, "POST", "/api/auth/register", map[string]string{
		"name":     "Refresh User",
		"email":    "refresh@example.com",
		"password": "Password123!",
	})

	var registerPayload authPayload
	require.NoError(t, json.Unmarshal(registerResp.Body.Bytes(), &registerPayload))

	resp := performJSONRequest(t, app, "POST", "/api/auth/refresh", map[string]string{
		"refresh_token": registerPayload.RefreshToken,
	})

	require.Equal(t, 200, resp.Code)

	var payload authPayload
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	assert.Equal(t, "success", payload.Status)
	assert.NotEmpty(t, payload.AccessToken)
	assert.NotEmpty(t, payload.RefreshToken)
	assert.Equal(t, "refresh@example.com", payload.User.Email)
}

func TestRefreshRejectsInvalidToken(t *testing.T) {
	app := setupAuthTestApp(t)

	resp := performJSONRequest(t, app, "POST", "/api/auth/refresh", map[string]string{
		"refresh_token": "invalid-token",
	})

	require.Equal(t, 401, resp.Code)

	var payload errorPayload
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	assert.Equal(t, "error", payload.Status)
}
