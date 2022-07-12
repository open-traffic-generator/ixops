package configs

func SetDefaultString(p **string, val string) {
	if *p == nil {
		*p = &val
	}
}

func SetDefaultBool(p **bool, val bool) {
	if *p == nil {
		*p = &val
	}
}

func SetDefaultInt(p **int, val int) {
	if *p == nil {
		*p = &val
	}
}
