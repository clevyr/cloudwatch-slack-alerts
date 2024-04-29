package alert

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/clevyr/cloudwatch-slack-alerts/internal/util"
	"github.com/slack-go/slack"
)

func (event Event) SlackMsg() []slack.MsgOption {
	var emoji string
	var titleExt string
	var color string
	switch event.AlarmData.State.Value {
	case StateAlarm:
		emoji = "red_circle"
		titleExt = "in alarm"
		color = "#df3617"
	case StateOK:
		emoji = "large_green_circle"
		titleExt = "OK"
		color = "#22af64"
	case StateInsufficientData:
		emoji = "large_yellow_circle"
		titleExt = "missing data"
		color = "#ffbf00"
	}

	consoleURL := url.URL{
		Scheme:   "https",
		Host:     "console.aws.amazon.com",
		Path:     "/cloudwatch/home",
		RawQuery: "region=" + event.Region,
		Fragment: "alarmsV2:alarm/" + event.AlarmData.AlarmName,
	}

	var context string
	if description := event.AlarmData.Configuration.Description; description != "" {
		context += "\n*Description:* " + description
	}
	switch event.AlarmData.State.Value {
	case StateAlarm, StateInsufficientData:
		context += "\n*Started at:* _" + event.AlarmData.State.Timestamp.Local().Format("3:04:05 PM") + "_"
	case StateOK:
		context += "\n*Duration:* " + event.AlarmData.State.Timestamp.Sub(event.AlarmData.PreviousState.Timestamp.Time).Round(time.Second).String()
	}
	if reason := event.AlarmData.State.Reason; reason != "" {
		context += "\n*Reason:* " + reason
	}
	if event.AlarmData.State.Value != StateOK && event.AlarmData.State.ReasonData != nil {
		if metrics := event.AlarmData.Configuration.Metrics; len(metrics) != 0 {
			datapoints := event.AlarmData.State.ReasonData.RecentDatapoints
			values := make([]string, 0, len(metrics))
			for i, metric := range metrics {
				if i < len(datapoints) {
					val := strconv.FormatFloat(datapoints[i], 'f', -1, 64)
					if name := metric.MetricStat.Metric.Name; len(metrics) > 1 || name != event.AlarmData.AlarmName {
						val = name + "=" + val
					}
					if unit := metric.MetricStat.Unit; unit != "" {
						val += " " + unit
					}
					values = append(values, val)
				}
			}
			if len(values) != 0 {
				label := util.Pluralize("Value", "Values", len(values))
				context += "\n*" + label + ":* " + strings.Join(values, ", ")
			}
		}
	}

	return []slack.MsgOption{
		slack.MsgOptionText("*AWS CloudWatch Notification*", false),
		slack.MsgOptionBlocks(
			slack.NewRichTextBlock("",
				slack.NewRichTextSection(
					slack.NewRichTextSectionEmojiElement(emoji, 2, nil),
					slack.NewRichTextSectionTextElement(
						" CloudWatch Metric "+titleExt+": ", &slack.RichTextSectionTextStyle{Bold: true},
					),
					slack.NewRichTextSectionTextElement(
						event.AlarmData.AlarmName, &slack.RichTextSectionTextStyle{Code: true},
					),
				),
			),
		),
		slack.MsgOptionAttachments(slack.Attachment{
			Color: color,
			Blocks: slack.Blocks{
				BlockSet: []slack.Block{
					slack.NewContextBlock("",
						slack.NewTextBlockObject(slack.MarkdownType, context, false, false),
					),
					slack.NewActionBlock("", &slack.ButtonBlockElement{
						Type:  slack.METButton,
						Text:  slack.NewTextBlockObject(slack.PlainTextType, "View Alarm", false, false),
						Style: slack.StylePrimary,
						URL:   consoleURL.String(),
					}),
				},
			},
		}),
	}
}
