package object

func IsErrorValue(obj Object) bool {
	if _, ok := obj.(*Error); ok {
		return true
	}
	return false
}

func IsNull(obj Object) bool {
	return obj == NULL
}

func IsReturnValue(obj Object) bool {
	if _, ok := obj.(*ReturnValue); ok {
		return true
	}
	return false
}

func IsTruthy(obj Object) bool {
	switch obj {
	case NULL, FALSE:
		return false
	default:
		return true
	}
}
