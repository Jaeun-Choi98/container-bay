package response

type BaseResponse struct {
	Result int `json:"result"`
	Data   any `json:"data"`
}

var (
	// user define stats
	SUCCESS       = NewBaseResponse(0)
	FAIL          = NewBaseResponse(1)
	NACK          = NewBaseResponse(2)
	TIMEOUT       = NewBaseResponse(3)
	TIMEOUT_WEB   = NewBaseResponse(4)
	INVAILD_DATA  = NewBaseResponse(5)
	INVALID_TOKEN = NewBaseResponse(6)

	INVAILD_DB_TYPE = NewBaseResponse(21)

	// 로그인 관련
	INVALID_LOGIN_ID  = NewBaseResponse(401).Add("ID가 존재하지 않습니다.")
	INVAILD_LOGIN_PWD = NewBaseResponse(401).Add("비밀번호가 일치하지 않습니다.")
	EXISTS_SESSION    = NewBaseResponse(303).Add("이미 로그인한 사용자입니다. 로그인 하시겠습니까?")
)

func NewBaseResponse(result int) *BaseResponse {
	return &BaseResponse{
		Result: result,
	}
}

func (r *BaseResponse) Add(data any) *BaseResponse {
	n := NewBaseResponse(r.Result)
	n.Data = data
	return n
}
