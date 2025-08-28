package cmd

import (
	"fmt"
	"os"
	"search-ads-reporter-ui/internal/config"
	"search-ads-reporter-ui/internal/sendgrid"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "reporter",
	Short: "A CLI for sending reports",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var sendTestEmailCmd = &cobra.Command{
	Use:   "send-test-email [to-email]",
	Short: "Send a test email",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig("configs/config.json")
		if err != nil {
			return fmt.Errorf("could not load config: %w", err)
		}

		if cfg.SendGridAPIKey == "YOUR_SENDGRID_API_KEY" {
			return fmt.Errorf("sendgrid api key is not configured in configs/config.json")
		}

		fromEmail := "reports@example.com"
		toEmail := args[0]
		subject := "Test Email"
		body := "<h1>This is a test email</h1>"

		return sendgrid.SendEmail(cfg.SendGridAPIKey, fromEmail, toEmail, subject, body)
	},
}

func init() {
	rootCmd.AddCommand(sendTestEmailCmd)
}
