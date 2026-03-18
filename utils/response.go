package utils

func NewResponse[T string | any](message, err string, data T) map[string]interface{} {
	if err != "" {
		return map[string]interface{}{
			"message": message,
			"error":   err,
		}
	}

	return map[string]interface{}{
		"message": message,
		"data":    data,
	}
}
