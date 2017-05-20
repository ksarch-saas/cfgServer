package api


type Response struct {
	Errno  int         `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Body   interface{} `json:"body"`
}

func MakeResponse(errno int, msg string, body interface{}) Response {
	return Response{errno, msg, body}
}

func MakeSuccessResponse(body interface{}) Response {
	return MakeResponse(0, "OK", body)
}

func MakeFailureResponse(msg string) Response {
	return MakeResponse(777, msg, nil)
}

type UpdateNodesParams struct {
	Region 			string       	`json:"region"`
	CfgID			string			`json:"CfgID"`
	Seeds  			interface{} 	`json:"seeds"`
}