package model

type CustomerField struct {
	Id       int    `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type ReturnUserData struct {
	Id       int    `json:"id,omitempty"`
	UserName string `json:"username,omitempty"`
	Role     string `json:"role,omitempty"`
	Token    string `json:"token,omitempty"`
}
