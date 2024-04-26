package alert

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/clevyr/cloudwatch-slack-alerts/internal/config"
	"github.com/slack-go/slack"
)

type HandlerFunc func(ctx context.Context, event events.SNSEvent) error

func Handler(conf *config.Config) HandlerFunc {
	var api *slack.Client
	return func(ctx context.Context, event events.SNSEvent) error {
		if api == nil {
			api = slack.New(conf.SlackAPIToken)
		}

		var errs []error
		for _, record := range event.Records {
			if err := notify(ctx, conf, api, record); err != nil {
				errs = append(errs, err)
			}
		}
		return errors.Join(errs...)
	}
}

func notify(ctx context.Context, conf *config.Config, api *slack.Client, record events.SNSEventRecord) error {
	opts, err := SlackMsg(record)
	if err != nil {
		return err
	}
	opts = append(opts, slack.MsgOptionAsUser(true))

	_, _, err = api.PostMessageContext(ctx, conf.SlackChannel, opts...)
	return err
}
