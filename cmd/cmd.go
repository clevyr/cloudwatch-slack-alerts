package cmd

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/clevyr/cloudwatch-slack-alerts/internal/alert"
	"github.com/clevyr/cloudwatch-slack-alerts/internal/config"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	conf := config.New()

	cmd := &cobra.Command{
		Use:   "cloudwatch-slack-alerts",
		Short: "Send AWS CloudWatch notifications to a Slack channel using Lambda",
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return conf.Load(cmd)
		},
		Run: func(_ *cobra.Command, _ []string) {
			lambda.Start(alert.Handler(conf))
		},
	}

	conf.RegisterFlags(cmd)
	cmd.SetContext(config.NewContext(context.Background(), conf))
	return cmd
}
