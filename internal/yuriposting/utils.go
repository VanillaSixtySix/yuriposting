package yuriposting

func Pluralize(num int) string {
	if num > 1 {
		return "s"
	}
	return ""
}
