package format_api

const CodeSuccess int64 = 10000
const CommonError int64 = 100000
const ParameterError int64 = 200000
const AuthorizationFail int64 = 250000
const ServerError int64 = 300000

func CodeString(code int64) string {
	switch code {
	case CodeSuccess:
		return "success"
	case CommonError:
		return "common error"
	default:
		return "system error"
	}
}
