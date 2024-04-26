package alert

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack"
)

func SlackMsg(snsEvent events.SNSEventRecord) ([]slack.MsgOption, error) {
	var event events.CloudWatchAlarmSNSPayload
	if err := json.Unmarshal([]byte(snsEvent.SNS.Message), &event); err != nil {
		return nil, err
	}

	var emoji string
	var titleExt string
	var color string
	switch event.NewStateValue {
	case "ALARM":
		emoji = "red_circle"
		titleExt = "in alarm"
		color = "#df3617"
	case "OK":
		emoji = "large_green_circle"
		titleExt = "OK"
		color = "#22af64"
	case "INSUFFICIENT_DATA":
		emoji = "large_yellow_circle"
		titleExt = "missing data"
		color = "#ffbf00"
	}

	consoleURL := url.URL{
		Scheme:   "https",
		Host:     "console.aws.amazon.com",
		Path:     "/cloudwatch/home",
		RawQuery: "region=" + strings.Split(snsEvent.EventSubscriptionArn, ":")[3],
		Fragment: "alarmsV2:alarm/" + event.AlarmName,
	}

	var context string
	if event.AlarmDescription != "" {
		context += "\n*Description:* " + event.AlarmDescription
	}
	context += "\n*Trigger:* " + formatStatistic(event.Trigger.Statistic) + " " +
		event.Trigger.MetricName + " " +
		formatComparison(event.Trigger.ComparisonOperator) + " " +
		strconv.FormatFloat(event.Trigger.Threshold, 'f', 0, 64) + " for "
	if event.Trigger.EvaluationPeriods > 1 {
		context += strconv.Itoa(int(event.Trigger.EvaluationPeriods)) + " periods of "
	}
	period := time.Duration(event.Trigger.Period) * time.Second
	context += strings.TrimSuffix(period.String(), "0s")
	if event.NewStateReason != "" {
		context += "\n*Reason:* " + event.NewStateReason
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
						event.AlarmName, &slack.RichTextSectionTextStyle{Code: true},
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
	}, nil
}

func formatComparison(comparison string) string {
	switch comparison {
	case "GreaterThanOrEqualToThreshold":
		return ">="
	case "GreaterThanThreshold":
		return ">"
	case "LessThanThreshold":
		return "<"
	case "LessThanOrEqualToThreshold":
		return "<="
	case "LessThanLowerOrGreaterThanUpperThreshold":
		return "outside threshold"
	case "LessThanLowerThreshold":
		return "below threshold"
	case "GreaterThanUpperThreshold":
		return "above threshold"
	default:
		return comparison
	}
}

func formatStatistic(statistic string) string {
	switch strings.ToLower(statistic) {
	case "samplecount":
		return "Sample count"
	case "average":
		return "Average"
	case "sum":
		return "Sum"
	case "minimum":
		return "Min"
	case "maximum":
		return "Max"
	default:
		return statistic
	}
}
