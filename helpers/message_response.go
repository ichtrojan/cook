package helpers

func PrepareMessage(message string) map[string]string {
	response := make(map[string]string)
	response["message"] = message
	return response
}

func PrepareResponse(key string, value interface{}) map[string]interface{} {
	response := make(map[string]interface{})
	response[key] = value
	return response
}
