package user

type UserInfo struct {
	Id         int64  `json:"id"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	Mobile     string `json:"mobile"`
	MobileCode string `json:"mobile_code"`
	Email      string `json:"email"`
	CountryId  int    `json:"country_id"`
	Country    string `json:"country"`
	CountryEn  string `json:"country_en"`
	ParentId   int64  `json:"parent_id"`
	ParentTree string `json:"parent_tree"`
	Level      int    `json:"level"`
	InviteCode string `json:"invite_code"`
	Status     uint   `json:"status"`
	AuthStatus int    `json:"auth_status"`
	HeadImg    string `json:"head_img"`
	UserLevel  int    `json:"user_level"`
}
