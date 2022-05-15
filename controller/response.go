package controller

import "TikTok/vo"

type UserLoginResponse struct {
	vo.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	vo.Response
	User vo.User `json:"user"`
}
