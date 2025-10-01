package error_common

var (
	SUCCESS    = 0
	ERR_COMMON = 400
)

var messages = map[int]map[string]string{
	SUCCESS: {
		"vi": "Thành công",
		"en": "Success",
	},
}
