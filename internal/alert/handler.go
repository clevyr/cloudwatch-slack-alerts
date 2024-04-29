package alert

import (
	"context"

	"github.com/clevyr/cloudwatch-slack-alerts/internal/config"
	"github.com/slack-go/slack"
)

type HandlerFunc func(ctx context.Context, event Event) error

func Handler(conf *config.Config) HandlerFunc {
	var api *slack.Client
	return func(ctx context.Context, event Event) error {
		if api == nil {
			api = slack.New(conf.SlackAPIToken)
		}

		opts := append(event.SlackMsg(), slack.MsgOptionAsUser(true))
		_, _, err := api.PostMessageContext(ctx, conf.SlackChannel, opts...)
		return err
	}
}
