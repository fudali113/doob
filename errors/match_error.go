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

func (mme *URLMacthError) Error() string {
	return mme.url
}

type MethodMacthError struct {
	shouldMethod string
	factMethod   string
	matchError   *URLMacthError
}

func (mme *MethodMacthError) Error() string {
	return mme.shouldMethod
}
