// Package usecases implements business logic for collections client-late-paying.
package usecases

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hawthorne/bff-api-go-collections-client-late-paying/internal/infrastructure"
)

// CollectionsUseCase handles all collection-related business logic.
type CollectionsUseCase struct {
	core *infrastructure.CoreClient
}

// NewCollectionsUseCase creates a new CollectionsUseCase.
func NewCollectionsUseCase(core *infrastructure.CoreClient) *CollectionsUseCase {
	return &CollectionsUseCase{core: core}
}

// ForwardActivity forwards an activity request to core-api.
func (uc *CollectionsUseCase) ForwardActivity(action string, activityID string, body interface{}) (interface{}, error) {
	path := "/api/v1/collections/activities"
	if activityID != "" {
		path = fmt.Sprintf("/api/v1/collections/activities/%s", activityID)
		if action == "patch" {
			path += "" // PATCH /activities/{id}
		}
	}

	var resp *http.Response
	var err error

	switch action {
	case "post":
		resp, err = uc.core.ForwardRequest("POST", path, body)
	case "patch":
		resp, err = uc.core.ForwardRequest("PATCH", path, body)
	case "get":
		resp, err = uc.core.ForwardRequest("GET", path, nil)
	default:
		return nil, fmt.Errorf("unsupported activity action: %s", action)
	}

	if err != nil {
		return nil, fmt.Errorf("forward activity %s: %w", action, err)
	}
	defer resp.Body.Close()

	return handleCoreResponse(resp)
}

// ForwardEscalation forwards an escalation request to core-api.
func (uc *CollectionsUseCase) ForwardEscalation(action string, escalationID string, subAction string, body interface{}) (interface{}, error) {
	path := "/api/v1/collections/escalations"
	if escalationID != "" {
		path = fmt.Sprintf("/api/v1/collections/escalations/%s", escalationID)
		if subAction != "" {
			path += "/" + subAction
		}
	}

	var resp *http.Response
	var err error

	method := "POST"
	if action == "get" {
		method = "GET"
	}

	resp, err = uc.core.ForwardRequest(method, path, body)
	if err != nil {
		return nil, fmt.Errorf("forward escalation %s: %w", action, err)
	}
	defer resp.Body.Close()

	return handleCoreResponse(resp)
}

// ForwardAccount forwards an account request to core-api.
func (uc *CollectionsUseCase) ForwardAccount(action string, accountID string, subAction string, body interface{}) (interface{}, error) {
	path := "/api/v1/collections/portfolio/accounts"
	if accountID != "" {
		path = fmt.Sprintf("/api/v1/collections/portfolio/accounts/%s", accountID)
		if subAction != "" {
			path += "/" + subAction
		}
	}

	var resp *http.Response
	var err error

	method := "GET"
	if action == "post" {
		method = "POST"
	}

	resp, err = uc.core.ForwardRequest(method, path, body)
	if err != nil {
		return nil, fmt.Errorf("forward account %s: %w", action, err)
	}
	defer resp.Body.Close()

	return handleCoreResponse(resp)
}

// ForwardReport forwards a report request to core-api.
func (uc *CollectionsUseCase) ForwardReport(reportType string, agentID string) (interface{}, error) {
	path := fmt.Sprintf("/api/v1/collections/reports/%s", reportType)
	if agentID != "" {
		path = fmt.Sprintf("/api/v1/collections/reports/%s/%s", reportType, agentID)
	}

	resp, err := uc.core.ForwardRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("forward report %s: %w", reportType, err)
	}
	defer resp.Body.Close()

	return handleCoreResponse(resp)
}

// ForwardUser forwards a user management request to core-api.
func (uc *CollectionsUseCase) ForwardUser(action string, userID string, body interface{}) (interface{}, error) {
	path := "/api/v1/admin/users"
	if userID != "" {
		path = fmt.Sprintf("/api/v1/admin/users/%s", userID)
	}

	var resp *http.Response
	var err error

	method := "GET"
	if action == "post" {
		method = "POST"
	} else if action == "patch" {
		method = "PATCH"
	}

	resp, err = uc.core.ForwardRequest(method, path, body)
	if err != nil {
		return nil, fmt.Errorf("forward user %s: %w", action, err)
	}
	defer resp.Body.Close()

	return handleCoreResponse(resp)
}

// ForwardConfig forwards a config request to core-api.
func (uc *CollectionsUseCase) ForwardConfig(action string, body interface{}) (interface{}, error) {
	path := "/api/v1/admin/config"

	var resp *http.Response
	var err error

	method := "GET"
	if action == "put" {
		method = "PUT"
	}

	resp, err = uc.core.ForwardRequest(method, path, body)
	if err != nil {
		return nil, fmt.Errorf("forward config %s: %w", action, err)
	}
	defer resp.Body.Close()

	return handleCoreResponse(resp)
}

// ForwardTenant forwards a tenant request to core-api.
func (uc *CollectionsUseCase) ForwardTenant(action string, body interface{}) (interface{}, error) {
	path := "/api/v1/admin/tenants"

	resp, err := uc.core.ForwardRequest(action, path, body)
	if err != nil {
		return nil, fmt.Errorf("forward tenant %s: %w", action, err)
	}
	defer resp.Body.Close()

	return handleCoreResponse(resp)
}

// handleCoreResponse reads and parses the core response.
func handleCoreResponse(resp *http.Response) (interface{}, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read core response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("core returned status %d: %s", resp.StatusCode, string(body))
	}

	var result interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("unmarshal core response: %w", err)
		}
	}
	return result, nil
}
