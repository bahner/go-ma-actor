package config

var template map[string]interface{}

func SetTemplate(t map[string]interface{}) {
	template = t
}

func GetTemplate() map[string]interface{} {
	return template
}
