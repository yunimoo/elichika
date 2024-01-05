package enum

const (
	MinuteSecondCount = 60
	HourMinuteCount   = 60
	HourSecondCount   = HourMinuteCount * MinuteSecondCount
	DayHourCount      = 24
	DaySecondCount    = DayHourCount * HourSecondCount
	WeekDayCount      = 7
	WeekSecondCount   = WeekDayCount * DaySecondCount
)
