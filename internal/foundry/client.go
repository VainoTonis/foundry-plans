package foundry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type PlanStep struct {
	ID       int64  `json:"id"`
	PlanID   int64  `json:"plan_id"`
	Position int    `json:"position"`
	Text     string `json:"text"`
	Status   string `json:"status"`
}

type Plan struct {
	ID       int64       `json:"id"`
	RepoName string      `json:"repo_name"`
	Title    string      `json:"title"`
	Summary  string      `json:"summary"`
	Status   string      `json:"status"`
	Steps    []PlanStep  `json:"steps"`
}

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		http:    &http.Client{},
	}
}

func (c *Client) ListPlans() ([]Plan, error) {
	resp, err := c.http.Get(c.baseURL + "/api/plans")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list plans failed: %d: %s", resp.StatusCode, string(body))
	}

	var plans []Plan
	if err := json.NewDecoder(resp.Body).Decode(&plans); err != nil {
		return nil, err
	}
	return plans, nil
}

func (c *Client) CreatePlan(repoName, title, summary string, steps []string) (*Plan, error) {
	payload := map[string]interface{}{
		"repo_name": repoName,
		"title":     title,
		"summary":   summary,
	}

	body, _ := json.Marshal(payload)
	resp, err := c.http.Post(
		c.baseURL+"/api/plans",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create plan failed: %d: %s", resp.StatusCode, string(respBody))
	}

	var plan Plan
	if err := json.NewDecoder(resp.Body).Decode(&plan); err != nil {
		return nil, err
	}

	// POST each step
	for _, stepText := range steps {
		stepPayload := map[string]interface{}{
			"text":   stepText,
			"status": "pending",
		}
		stepBody, _ := json.Marshal(stepPayload)
		stepResp, err := c.http.Post(
			fmt.Sprintf("%s/api/plans/%d/steps", c.baseURL, plan.ID),
			"application/json",
			bytes.NewReader(stepBody),
		)
		if err != nil {
			return nil, err
		}
		if stepResp.StatusCode != http.StatusCreated && stepResp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(stepResp.Body)
			stepResp.Body.Close()
			return nil, fmt.Errorf("create step failed: %d: %s", stepResp.StatusCode, string(respBody))
		}
		stepResp.Body.Close()
	}

	// GET plan with steps
	return c.GetPlan(plan.ID)
}

func (c *Client) GetPlan(id int64) (*Plan, error) {
	resp, err := c.http.Get(fmt.Sprintf("%s/api/plans/%d", c.baseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get plan failed: %d: %s", resp.StatusCode, string(body))
	}

	var plan Plan
	if err := json.NewDecoder(resp.Body).Decode(&plan); err != nil {
		return nil, err
	}

	// GET steps
	stepsResp, err := c.http.Get(fmt.Sprintf("%s/api/plans/%d/steps", c.baseURL, id))
	if err != nil {
		return nil, err
	}
	defer stepsResp.Body.Close()

	if stepsResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(stepsResp.Body)
		return nil, fmt.Errorf("get steps failed: %d: %s", stepsResp.StatusCode, string(body))
	}

	var steps []PlanStep
	if err := json.NewDecoder(stepsResp.Body).Decode(&steps); err != nil {
		return nil, err
	}
	plan.Steps = steps

	return &plan, nil
}

func (c *Client) UpdatePlan(id int64, updates map[string]interface{}) (*Plan, error) {
	body, _ := json.Marshal(updates)

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("%s/api/plans/%d", c.baseURL, id), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("update plan failed: %d: %s", resp.StatusCode, string(respBody))
	}

	var plan Plan
	if err := json.NewDecoder(resp.Body).Decode(&plan); err != nil {
		return nil, err
	}
	return &plan, nil
}

func (c *Client) UpdateStep(planID, stepID int64, status, text string) (*PlanStep, error) {
	payload := map[string]interface{}{}
	if status != "" {
		payload["status"] = status
	}
	if text != "" {
		payload["text"] = text
	}

	return c.UpdateStepFromMap(planID, stepID, payload)
}

func (c *Client) UpdateStepFromMap(planID, stepID int64, updates map[string]interface{}) (*PlanStep, error) {
	body, _ := json.Marshal(updates)

	endpoint := fmt.Sprintf("%s/api/plans/%d/steps/%d", c.baseURL, planID, stepID)
	req, _ := http.NewRequest("PATCH", endpoint, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("update step failed: %d: %s", resp.StatusCode, string(respBody))
	}

	var step PlanStep
	if err := json.NewDecoder(resp.Body).Decode(&step); err != nil {
		return nil, err
	}
	return &step, nil
}
