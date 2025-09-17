package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootRedirect(t *testing.T) {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/web/home")
	})
	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Serve the HTTP request
	r.ServeHTTP(rr, req)
	// Check the status code for redirect
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}
	// Check the redirect location
	location := rr.Header().Get("Location")
	expected := "/web/home"
	if location != expected {
		t.Errorf("Handler returned unexpected redirect location: got %v want %v", location, expected)
	}
}
