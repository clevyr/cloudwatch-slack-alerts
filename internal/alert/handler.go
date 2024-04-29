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

		if err := notify(ctx, conf, api, event); err != nil {
			return err
		}
		return nil
	}
}

func notify(ctx context.Context, conf *config.Config, api *slack.Client, event Event) error {
	opts, err := event.SlackMsg()
	if err != nil {
		return err
	}
	opts = append(opts, slack.MsgOptionAsUser(true))

	_, _, err = api.PostMessageContext(ctx, conf.SlackChannel, opts...)
	return err
}
