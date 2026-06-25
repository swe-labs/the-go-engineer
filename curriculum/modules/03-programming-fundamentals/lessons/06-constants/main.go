package main

import "fmt"

const (
	DaysInWeek      = 7
	HoursInDay      = 24
	MinutesInHour   = 60
	SecondsInMinute = 60
	SecondsPerWeek  = DaysInWeek * HoursInDay * MinutesInHour * SecondsInMinute
)

const (
	Monday = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func main() {
	fmt.Println("Seconds per week:", SecondsPerWeek)

	fmt.Printf("Monday = %d\n", Monday)
	fmt.Printf("Tuesday = %d\n", Tuesday)
	fmt.Printf("Wednesday = %d\n", Wednesday)
	fmt.Printf("Thursday = %d\n", Thursday)
	fmt.Printf("Friday = %d\n", Friday)
	fmt.Printf("Saturday = %d\n", Saturday)
	fmt.Printf("Sunday = %d\n", Sunday)
}

func weekdayName(n int) string {
	switch n {
	case Monday:
		return "Monday"
	case Tuesday:
		return "Tuesday"
	case Wednesday:
		return "Wednesday"
	case Thursday:
		return "Thursday"
	case Friday:
		return "Friday"
	case Saturday:
		return "Saturday"
	case Sunday:
		return "Sunday"
	default:
		return "unknown"
	}
}
