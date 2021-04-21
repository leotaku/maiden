package timeutil

type DateOrder func(string, string, string) (string, string, string)

func ISO(year, month, day string) (string, string, string) {
	return year, month, day
}

func European(day, month, year string) (string, string, string) {
	return year, month, day
}

func American(month, day, year string) (string, string, string) {
	return year, month, day
}

