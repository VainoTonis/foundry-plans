package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/VainoTonis/foundry-plans/internal/foundry"
	"github.com/spf13/cobra"
)

var plansCmd = &cobra.Command{
	Use:   "plans",
	Short: "Manage plans",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all plans",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := foundry.NewClient(apiURL)
		plans, err := client.ListPlans()
		if err != nil {
			return err
		}
		output, _ := json.MarshalIndent(plans, "", "  ")
		fmt.Println(string(output))
		return nil
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new plan (JSON from stdin only)",
	Long:  "Create a new plan from JSON input on stdin.\n\nRequired JSON fields: repo_name (string), title (string)\nOptional JSON fields: summary (string), steps (array of strings)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %v", err)
		}

		var req map[string]interface{}
		if err := json.Unmarshal(input, &req); err != nil {
			return fmt.Errorf("invalid JSON input: %v", err)
		}

		repoName, ok := req["repo_name"].(string)
		if !ok || repoName == "" {
			return fmt.Errorf("repo_name is required")
		}

		title, ok := req["title"].(string)
		if !ok || title == "" {
			return fmt.Errorf("title is required")
		}

		summary, _ := req["summary"].(string)

		var steps []string
		if stepsInterface, ok := req["steps"]; ok {
			if stepsArray, ok := stepsInterface.([]interface{}); ok {
				for _, step := range stepsArray {
					if stepStr, ok := step.(string); ok {
						steps = append(steps, stepStr)
					}
				}
			}
		}

		client := foundry.NewClient(apiURL)
		plan, err := client.CreatePlan(repoName, title, summary, steps)
		if err != nil {
			return err
		}
		output, _ := json.MarshalIndent(plan, "", "  ")
		fmt.Println(string(output))
		return nil
	},
}

var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a plan by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var id int64
		if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil {
			return fmt.Errorf("invalid plan id: %v", err)
		}
		client := foundry.NewClient(apiURL)
		plan, err := client.GetPlan(id)
		if err != nil {
			return err
		}
		output, _ := json.MarshalIndent(plan, "", "  ")
		fmt.Println(string(output))
		return nil
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a plan (JSON from stdin only)",
	Long:  "Update a plan from JSON input on stdin.\n\nRequired JSON field: id (number)\nOptional JSON fields: status (string), title (string), summary (string)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %v", err)
		}

		var updates map[string]interface{}
		if err := json.Unmarshal(input, &updates); err != nil {
			return fmt.Errorf("invalid JSON input: %v", err)
		}

		idInterface, ok := updates["id"]
		if !ok {
			return fmt.Errorf("id is required")
		}

		var id int64
		switch v := idInterface.(type) {
		case float64:
			id = int64(v)
		case string:
			if _, err := fmt.Sscanf(v, "%d", &id); err != nil {
				return fmt.Errorf("invalid id: %v", err)
			}
		default:
			return fmt.Errorf("id must be a number or string")
		}

		// Remove id from updates before sending to API
		delete(updates, "id")

		if len(updates) == 0 {
			return fmt.Errorf("no fields to update provided")
		}

		client := foundry.NewClient(apiURL)
		plan, err := client.UpdatePlan(id, updates)
		if err != nil {
			return err
		}
		output, _ := json.MarshalIndent(plan, "", "  ")
		fmt.Println(string(output))
		return nil
	},
}

var updateStepCmd = &cobra.Command{
	Use:   "update-step",
	Short: "Update a step (JSON from stdin only)",
	Long:  "Update a step from JSON input on stdin.\n\nRequired JSON fields: plan_id (number), step_id (number)\nOptional JSON fields: status (string), text (string)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %v", err)
		}

		var updates map[string]interface{}
		if err := json.Unmarshal(input, &updates); err != nil {
			return fmt.Errorf("invalid JSON input: %v", err)
		}

		planIDInterface, ok := updates["plan_id"]
		if !ok {
			return fmt.Errorf("plan_id is required")
		}

		stepIDInterface, ok := updates["step_id"]
		if !ok {
			return fmt.Errorf("step_id is required")
		}

		var planID, stepID int64
		
		switch v := planIDInterface.(type) {
		case float64:
			planID = int64(v)
		case string:
			if _, err := fmt.Sscanf(v, "%d", &planID); err != nil {
				return fmt.Errorf("invalid plan_id: %v", err)
			}
		default:
			return fmt.Errorf("plan_id must be a number or string")
		}

		switch v := stepIDInterface.(type) {
		case float64:
			stepID = int64(v)
		case string:
			if _, err := fmt.Sscanf(v, "%d", &stepID); err != nil {
				return fmt.Errorf("invalid step_id: %v", err)
			}
		default:
			return fmt.Errorf("step_id must be a number or string")
		}

		// Remove plan_id and step_id from updates before sending to API
		delete(updates, "plan_id")
		delete(updates, "step_id")

		if len(updates) == 0 {
			return fmt.Errorf("no fields to update provided")
		}

		client := foundry.NewClient(apiURL)
		step, err := client.UpdateStepFromMap(planID, stepID, updates)
		if err != nil {
			return err
		}
		output, _ := json.MarshalIndent(step, "", "  ")
		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(plansCmd)
	plansCmd.AddCommand(listCmd)
	plansCmd.AddCommand(createCmd)
	plansCmd.AddCommand(getCmd)
	plansCmd.AddCommand(updateCmd)
	plansCmd.AddCommand(updateStepCmd)
}
