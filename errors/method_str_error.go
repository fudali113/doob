package errors

type MethodStrBarbarismError struct {
}

func (this *MethodStrBarbarismError) Error() string {
	return "barbarism http method"
}
