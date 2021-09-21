package response

var MsgMap = map[int]string{
	SUCCESS: "OK",
	FAIL:    "FAIL",
}

func GetMsg(code int) string {
	msg, ok := MsgMap[code]

	if ok {
		return msg
	}

	return MsgMap[FAIL]
}
