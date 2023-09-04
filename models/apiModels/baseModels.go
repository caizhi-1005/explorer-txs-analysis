package apiModels

//定义JSON返回格式
type DataResponse struct {
	Code  int         `json:"code"`
	EnMsg string      `json:"enMsg"`
	Data  interface{} `json:"data"`
}

//定义服务层返回信息
type DataResponseSer struct {
	Err     error       `json:"err"`
	ErrCode int         `json:"errCode"`
	EnMsg   string      `json:"enMsg"`
	Results interface{} `json:"results"`
}

//定义APIs返回信息
type APIsResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

//定义APIs返回信息
type JsonRpcResponse struct {
	ID      int         `json:"id"`
	Version string      `json:"jsonrpc"`
	Error   interface{} `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

type ReqCommon struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type ReqPagination struct {
	Start    string `json:"start"`
	Length   string `json:"length"`
	Page     int64  `json:"page"`
	PageSize int64  `json:"page_size"`
}
