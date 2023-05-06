package handlers

import (
	"config-generator/models"
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConfigurationHandler(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{
		"appName": "myApp",
		"version": "v1",
		"clusterUrl": "https://mycluster.com",
		"microservices": {
			"service1": {
				"serviceName": "myService1",
				"image": "myImage1"
			},
			"service2": {
				"serviceName": "myService2",
				"image": "myImage2"
			}
		}
	}`))

	// Set the request's content type to JSON
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Create a new HTTP response recorder
	rec := httptest.NewRecorder()

	// Create an Echo context with the request and response
	c := e.NewContext(req, rec)

	// Call the function to be tested
	err := ConfigurationHandler(c)

	// Check that the function returned no error
	if err != nil {
		t.Errorf("ConfigurationHandler returned an error: %v", err)
	}

	// Check the HTTP status code of the response
	if rec.Code != http.StatusOK {
		t.Errorf("Expected HTTP status code %d but got %d", http.StatusOK, rec.Code)
	}

	// Decode the response body into a RepoResponse struct
	var resp models.RepoResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// Check the HelmRepo and ArgoRepo fields of the response
	expectedHelmRepo := "https://github.com/Shenali-SJ/-helm.git"
	expectedArgoRepo := "https://github.com/Shenali-SJ/-argo.git"
	if resp.HelmRepo != expectedHelmRepo {
		t.Errorf("Expected HelmRepo to be %s but got %s", expectedHelmRepo, resp.HelmRepo)
	}
	if resp.ArgoRepo != expectedArgoRepo {
		t.Errorf("Expected ArgoRepo to be %s but got %s", expectedArgoRepo, resp.ArgoRepo)
	}
}
