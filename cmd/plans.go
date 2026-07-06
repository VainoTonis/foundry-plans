package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

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
	Short: "Create a new plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoName, _ := cmd.Flags().GetString("repo-name")
		title, _ := cmd.Flags().GetString("title")
		summary, _ := cmd.Flags().GetString("summary")
		stepsStr, _ := cmd.Flags().GetString("steps")

		if repoName == "" || title == "" {
			return fmt.Errorf("--repo-name and --title are required")
		}

		var steps []string
		if stepsStr != "" {
			steps = strings.Split(stepsStr, ",")
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

var updateStatusCmd = &cobra.Command{
	Use:   "update-status <id> <status>",
	Short: "Update plan status",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var id int64
		if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil {
			return fmt.Errorf("invalid plan id: %v", err)
		}
		client := foundry.NewClient(apiURL)
		plan, err := client.UpdatePlanStatus(id, args[1])
		if err != nil {
			return err
		}
		output, _ := json.MarshalIndent(plan, "", "  ")
		fmt.Println(string(output))
		return nil
	},
}

var updateStepCmd = &cobra.Command{
	Use:   "update-step <plan-id> <step-id> <status> [text]",
	Short: "Update step status and/or text",
	Args:  cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		var planID, stepID int64
		if _, err := fmt.Sscanf(args[0], "%d", &planID); err != nil {
			return fmt.Errorf("invalid plan id: %v", err)
		}
		if _, err := fmt.Sscanf(args[1], "%d", &stepID); err != nil {
			return fmt.Errorf("invalid step id: %v", err)
		}
		status := args[2]
		text := ""
		if len(args) > 3 {
			text = args[3]
		}

		client := foundry.NewClient(apiURL)
		step, err := client.UpdateStep(planID, stepID, status, text)
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
	plansCmd.AddCommand(updateStatusCmd)
	plansCmd.AddCommand(updateStepCmd)

	createCmd.Flags().String("repo-name", "", "Repository name (required)")
	createCmd.Flags().String("title", "", "Plan title (required)")
	createCmd.Flags().String("summary", "", "Plan summary")
	createCmd.Flags().String("steps", "", "Comma-separated list of steps")
}
