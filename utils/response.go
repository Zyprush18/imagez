package utils


func NewResponse[T string | any ](message string, data T) map[string]interface{} {
	return map[string]interface{}{
		"message": message,
		"data": data,
	}
}