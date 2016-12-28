package errors

type MatchError struct {
	message string
}

func (me *MatchError) Error() string {
	return me.message
}

type URLMacthError struct {
	url        string
	matchError *MatchError
}

func (me *URLMacthError) Error() string {
	return me.url
}

type MethodMacthError struct {
	shouldMethod string
	factMethod   string
	matchError   *URLMacthError
}

func (me *MethodMacthError) Error() string {
	return me.shouldMethod
}
