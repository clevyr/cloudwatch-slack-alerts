package alert

import (
	"encoding/json"
	"fmt"
)

type Event struct {
	Source    string    `json:"source"`
	AlarmARN  string    `json:"alarmArn"`
	AccountID string    `json:"accountId"`
	Time      Time      `json:"time"`
	Region    string    `json:"region"`
	AlarmData AlarmData `json:"alarmData"`
}

type AlarmData struct {
	AlarmName     string        `json:"alarmName"`
	State         State         `json:"state"`
	PreviousState State         `json:"previousState"`
	Configuration Configuration `json:"configuration"`
}

type State struct {
	Value      AlarmState  `json:"value"`
	Reason     string      `json:"reason"`
	ReasonData *ReasonData `json:"reasonData"`
	Timestamp  Time        `json:"timestamp"`
}

type ReasonData struct {
	Version             string               `json:"version"`
	QueryDate           Time                 `json:"queryDate"`
	StartDate           Time                 `json:"startDate"`
	Statistic           string               `json:"statistic"`
	Period              int                  `json:"period"`
	RecentDatapoints    []float64            `json:"recentDatapoints"`
	Threshold           float64              `json:"threshold"`
	EvaluatedDatapoints []EvaluatedDatapoint `json:"evaluatedDatapoints"`
}

func (r *ReasonData) UnmarshalJSON(b []byte) error {
	var data string
	if err := json.Unmarshal(b, &data); err != nil {
		return fmt.Errorf("failed to unmarshal reasonData from string: %w", err)
	}

	// Create temporary alias to prevent stack overflow
	type reasonData ReasonData
	if err := json.Unmarshal([]byte(data), (*reasonData)(r)); err != nil {
		return fmt.Errorf("failed to unmarshal reasonData to struct: %w", err)
	}

	return nil
}

type EvaluatedDatapoint struct {
	Timestamp   Time    `json:"timestamp"`
	SampleCount float64 `json:"sampleCount"`
	Value       float64 `json:"value"`
}

type Metric struct {
	Namespace  string         `json:"namespace"`
	Name       string         `json:"name"`
	Dimensions map[string]any `json:"dimensions"`
}

type MetricStat struct {
	Metric Metric `json:"metric"`
	Period int    `json:"period"`
	Stat   string `json:"stat"`
	Unit   string `json:"unit"`
}

type Metrics struct {
	ID         string     `json:"id"`
	MetricStat MetricStat `json:"metricStat"`
	ReturnData bool       `json:"returnData"`
}

type Configuration struct {
	Description string    `json:"description"`
	Metrics     []Metrics `json:"metrics"`
}
