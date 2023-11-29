package users

type ListRequest struct {
	PageSize  int64 `form:"page_size"`
	PageToken int64 `form:"page_token"`
}
