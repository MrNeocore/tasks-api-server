package util

func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetOrElse(f func(key string) (value string, found bool), key string, orElse string) string {
	if value, found := f(key); found {
		return value
	} else {
		return orElse
	}
}
