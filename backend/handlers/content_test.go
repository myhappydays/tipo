package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"tipo-backend/models"
)

func TestContentHandler(t *testing.T) {
	// Load .env for tests
	err := godotenv.Load("../../backend/.env") // Adjust path as necessary
	if err != nil {
		t.Log("Error loading .env file for tests, assuming env vars are set: ", err)
	}

	// Test Case 1: Successful request for book category
	t.Run("Successful Book Request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/content?category=book&query=%EC%96%B4%EB%A6%B0%EC%99%95%EC%9E%90", nil)
		rec := httptest.NewRecorder()

		ContentHandler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response models.ContentResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Title)
		assert.Equal(t, "도서", response.Category)
		assert.NotEmpty(t, response.Sentences)
	})

	// Test Case 2: Missing category parameter
	t.Run("Missing Category Parameter", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/content?query=test", nil)
		rec := httptest.NewRecorder()

		ContentHandler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "'category' and 'query' parameters are required")
	})

	// Test Case 3: Missing query parameter
	t.Run("Missing Query Parameter", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/content?category=book", nil)
		rec := httptest.NewRecorder()

		ContentHandler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "'category' and 'query' parameters are required")
	})

	// Test Case 4: Invalid category
	t.Run("Invalid Category", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/content?category=invalid&query=test", nil)
		rec := httptest.NewRecorder()

		ContentHandler(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid category")
	})

	

	// Test Case 6: Successful request for news category
	t.Run("Successful News Request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/content?category=news&query=%EC%84%B1%EB%8A%A5%20%ED%96%A5%EC%83%81", nil)
		rec := httptest.NewRecorder()

		ContentHandler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		var response models.ContentResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Title)
		assert.Equal(t, "뉴스", response.Category)
		assert.NotEmpty(t, response.Sentences)
	})
}
