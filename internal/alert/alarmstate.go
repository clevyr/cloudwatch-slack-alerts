package alert

type AlarmState string

const (
	StateOK               AlarmState = "OK"
	StateAlarm            AlarmState = "ALARM"
	StateInsufficientData AlarmState = "INSUFFICIENT_DATA"
)
