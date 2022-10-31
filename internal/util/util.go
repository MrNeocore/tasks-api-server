package util

import "strconv"

func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

type GetValueIsFoundFunc func(key string) (value string, found bool)

func GetOrElse(f GetValueIsFoundFunc, key string, orElse string) string {
	if value, found := f(key); found {
		return value
	} else {
		return orElse
	}
}

func StringToInt(s string) int {
	i, err := strconv.Atoi(s)

	if err != nil {
		PanicError(err)
	}

	return i
}
