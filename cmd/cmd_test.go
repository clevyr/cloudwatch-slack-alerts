package cmd

import (
	"testing"

	"github.com/clevyr/cloudwatch-slack-alerts/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCmdDbFlags(t *testing.T) {
	t.Setenv("SLACK_API_TOKEN", "123")
	t.Setenv("SLACK_CHANNEL", "alerts")

	rootCmd := NewCommand()
	rootCmd.Run = func(_ *cobra.Command, _ []string) {}

	require.NoError(t, rootCmd.Execute())

	conf, ok := config.FromContext(rootCmd.Context())
	require.True(t, ok)
	assert.Equal(t, "123", conf.SlackAPIToken)
	assert.Equal(t, "alerts", conf.SlackChannel)
}
