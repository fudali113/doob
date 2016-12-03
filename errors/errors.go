package errors

import "fmt"

type Error int

func (err Error) Error() string {
	return fmt.Sprintf("weight value is %d", int(err))
}

func (err Error) GetWV() int {
	return int(err)
}

func GetMwthodMatchError(should, fact, url string, errors ...error) *MethodMacthError {
	return &MethodMacthError{
		shouldMethod: should,
		factMethod:   fact,
		matchError: &URLMacthError{
			url: url,
			matchError: &MatchError{
				message: "",
			},
		},
	}
}

func GetURLMatchError(url, message string, errors ...error) *URLMacthError {
	return &URLMacthError{
		url: url,
		matchError: &MatchError{
			message: message,
		},
	}
}
