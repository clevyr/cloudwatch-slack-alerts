package alert

type AlarmState string

const (
	StateOK               AlarmState = "OK"
	StateAlarm            AlarmState = "ALARM"
	StateInsufficientData AlarmState = "INSUFFICIENT_DATA"
)

func (a AlarmState) TitleExt() string {
	switch a {
	case StateAlarm:
		return "in alarm"
	case StateOK:
		return "OK"
	case StateInsufficientData:
		return "missing data"
	default:
		return "state unknown"
	}
}

func (a AlarmState) Emoji() string {
	switch a {
	case StateAlarm:
		return "red_circle"
	case StateOK:
		return "large_green_circle"
	case StateInsufficientData:
		return "large_yellow_circle"
	default:
		return "white_circle"
	}
}

func (a AlarmState) Color() string {
	switch a {
	case StateAlarm:
		return "#df3617"
	case StateOK:
		return "#22af64"
	case StateInsufficientData:
		return "#ffbf00"
	default:
		return "#efefef"
	}
}
