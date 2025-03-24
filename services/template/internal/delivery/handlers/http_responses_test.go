package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpError(t *testing.T) {
	rr := httptest.NewRecorder()
	httpError(rr, http.StatusBadRequest, "Bad Request")

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := map[string]string{"error": "Bad Request"}
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if response["error"] != expected["error"] {
		t.Errorf("handler returned unexpected body: got %v want %v", response["error"], expected["error"])
	}
}

func TestHttpSuccess(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"message": "Success"}
	httpSuccess(rr, http.StatusOK, data)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if response["message"] != data["message"] {
		t.Errorf("handler returned unexpected body: got %v want %v", response["message"], data["message"])
	}
}
