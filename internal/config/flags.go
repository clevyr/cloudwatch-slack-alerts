package config

import "github.com/spf13/cobra"

const (
	SlackAPITokenFlag = "slack-api-token"
	SlackChannelFlag  = "slack-channel"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.SlackAPIToken, SlackAPITokenFlag, "", "Slack API token")
	cmd.Flags().StringVar(&c.SlackChannel, SlackChannelFlag, "", "Slack channel")
}
