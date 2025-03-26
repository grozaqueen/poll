package utils

import (
	"fmt"
	"strconv"
	"time"
)

var mskLocation = time.FixedZone("MSK", 3*60*60)

func ParseSimpleDate(dateStr string) (time.Time, error) {
	if len(dateStr) != 10 || dateStr[2] != '.' || dateStr[5] != '.' {
		return time.Time{}, fmt.Errorf("неверный формат даты, используйте DD.MM.YYYY")
	}

	day, err1 := strconv.Atoi(dateStr[:2])
	month, err2 := strconv.Atoi(dateStr[3:5])
	year, err3 := strconv.Atoi(dateStr[6:])

	if err1 != nil || err2 != nil || err3 != nil {
		return time.Time{}, fmt.Errorf("дата должна содержать только цифры и точки")
	}

	if day < 1 || day > 31 || month < 1 || month > 12 || year < 2000 {
		return time.Time{}, fmt.Errorf("некорректная дата")
	}

	return time.Date(
		year,
		time.Month(month),
		day,
		23, 59, 59, 0,
		mskLocation,
	), nil
}
