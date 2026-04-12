package bot

import (
	"strconv"
	"strings"
	"time"
)

// AvailableHours — рабочие часы для консультаций
var AvailableHours = []int{10, 11, 12, 14, 15, 16, 17}

// Slot представляет доступный слот для записи
type Slot struct {
	Time     time.Time
	Label    string
	Callback string
}

func GetAvailableSlots(lawyerID int, daysAhead int) []Slot {
	var slots []Slot
	now := time.Now()

	for day := 0; day < daysAhead; day++ {
		date := now.AddDate(0, 0, day)

		var dateLabel string
		if day == 0 {
			dateLabel = "Сегодня, " + date.Format("02.01")
		} else if day == 1 {
			dateLabel = "Завтра, " + date.Format("02.01")
		} else {
			dateLabel = date.Format("02.01.2006")
		}

		for _, hour := range AvailableHours {
			slotTime := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, time.Local)

			if slotTime.After(now) {
				slots = append(slots, Slot{
					Time:     slotTime,
					Label:    dateLabel + " " + formatHour(hour),
					Callback: formatCallback(lawyerID, slotTime),
				})
			}
		}
	}

	return slots
}

func formatHour(hour int) string {
	if hour < 10 {
		return "0" + itoa(hour) + ":00"
	}
	return itoa(hour) + ":00"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	tmp := n
	for tmp > 0 {
		digits = append([]byte{byte('0' + tmp%10)}, digits...)
		tmp /= 10
	}
	return string(digits)
}

func formatCallback(lawyerID int, t time.Time) string {
	return "slot_" + itoa(lawyerID) + "_" +
		itoa(t.Year()) + "_" +
		itoa(int(t.Month())) + "_" +
		itoa(t.Day()) + "_" +
		itoa(t.Hour())
}

func ParseCallback(callback string) (lawyerID int, t time.Time, ok bool) {
	parts := strings.Split(callback, "_")
	if len(parts) != 6 || parts[0] != "slot" {
		return 0, time.Time{}, false
	}

	var err error
	lawyerID, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, time.Time{}, false
	}
	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, time.Time{}, false
	}
	month, err := strconv.Atoi(parts[3])
	if err != nil {
		return 0, time.Time{}, false
	}
	day, err := strconv.Atoi(parts[4])
	if err != nil {
		return 0, time.Time{}, false
	}
	hour, err := strconv.Atoi(parts[5])
	if err != nil {
		return 0, time.Time{}, false
	}

	t = time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.Local)
	return lawyerID, t, true
}
