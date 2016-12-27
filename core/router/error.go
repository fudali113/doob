package router

type NotMatch struct {
	errorStr string
}

func (this NotMatch) Error() string {
	return this.errorStr
}
