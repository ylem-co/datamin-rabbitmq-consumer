package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type WorkflowClient interface {
	RunWorkflow(workflowUuid string, payload []byte) (string, error) // returns workflow run UUID
}

type workflowClient struct {
	apiClient   *resty.Client
	oauthClient OAuthClient
}

func (c *workflowClient) RunWorkflow(workflowUuid string, payload []byte) (string, error) {
	log.Tracef("Running a pipeline, UUID: %s", workflowUuid)
	t, err := c.oauthClient.GetToken()
	if err != nil {
		return "", err
	}

	log.Trace("Succcessfully obtained OAuth token")

	url := fmt.Sprintf("/v1/pipelines/%s/runs/", workflowUuid)
	resp, err := c.apiClient.
		R().
		SetAuthToken(t.AccessToken).
		SetBody(payload).
		SetHeader("Content-Type", "application/json").
		Post(url)

	if err != nil {
		return "", err
	}

	log.Trace("Succcessfully sent HTTP request")

	if resp.StatusCode() != http.StatusCreated {
		log.Tracef("Response: %s", string(resp.Body()))
		return "", fmt.Errorf("unexpected http status code: %d", resp.StatusCode())
	}

	log.Trace("Pipeline run initiated successfully")

	result := struct {
		WorkflowRunUuid string `json:"workflow_run_uuid"`
	}{}

	err = json.Unmarshal(resp.Body(), &result)

	return result.WorkflowRunUuid, err
}

func NewWorkflowClient(baseUrl, clientId, clientSecret, basicAuthUsername, basicAuthPassword string) WorkflowClient {
	log.Tracef("Creating a new API client, base url: %s", baseUrl)

	apiClient := resty.New().SetBaseURL(baseUrl)
	if len(basicAuthUsername) > 0 || len(basicAuthPassword) > 0 {
		apiClient.SetBasicAuth(basicAuthUsername, basicAuthPassword)
	}

	c := &workflowClient{
		apiClient:   apiClient,
		oauthClient: NewOAuthClient(baseUrl, clientId, clientSecret, basicAuthUsername, basicAuthPassword),
	}

	return c
}
