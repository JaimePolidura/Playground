package utils

type LoxError struct {
	Line    int
	Where   string
	Message string
}

func (e LoxError) Error() string {
	return ""
}
