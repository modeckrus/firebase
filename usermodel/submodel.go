package usermodel

//SubModel is model representing subscriber of user
type SubModel struct {
	UID    string `json:"uid"`
	Nick   string `json:"nick"`
	Avatar string `json:"avatar"`
}
