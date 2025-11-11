package common

import (
	"errors"
	"fmt"
)

func GetMonthNumber(month string) (string, error) {
	switch month {
	case "Jan", "January":
		return "01", nil
	case "Feb", "February":
		return "02", nil
	case "Mar", "March":
		return "03", nil
	case "Apr", "April":
		return "04", nil
	case "May":
		return "05", nil
	case "Jun", "June":
		return "06", nil
	case "Jul", "July":
		return "07", nil
	case "Aug", "August":
		return "08", nil
	case "Sep", "September":
		return "09", nil
	case "Oct", "October":
		return "10", nil
	case "Nov", "November":
		return "11", nil
	case "Dec", "December":
		return "12", nil
	default:
		msg := fmt.Sprintf("Unknown month: %s", month)
		return "", errors.New(msg)

	}
}
