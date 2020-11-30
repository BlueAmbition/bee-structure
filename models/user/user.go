package user

import (
	"bee-structure/functions/hash"
	"bee-structure/models/tool"
	"fmt"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
)

//用户 model
type User struct {
	Id          int64     `orm:"column(id)"`
	Avatar      string    `orm:"column(avatar)"`
	Nickname    string    `orm:"column(nickname)"`
	Mobile      string    `orm:"column(mobile)"`
	MobileCode  string    `orm:"column(mobile_code)"`
	Email       string    `orm:"column(email)"`
	CountryId   int64     `orm:"column(country_id)"`
	Password    string    `orm:"column(password)"`
	PayPassword string    `orm:"column(pay_password)"`
	ParentId    int64     `orm:"column(parent_id)"`
	ParentTree  string    `orm:"column(parent_tree)"`
	Level       int       `orm:"column(level)"`
	InviteCode  string    `orm:"column(invite_code)"`
	Status      uint      `orm:"column(status)"`
	HeadImg     string    `orm:"column(head_img)"`
	IMEI        string    `orm:"column(imei)"`
	CreatedAt   time.Time `orm:"column(created_at);auto_now_add;type(datetime)"`
	UpdatedAt   time.Time `orm:"column(updated_at);auto_now;type(datetime)"`
	InvitedAt   time.Time `orm:"column(invited_at);auto_now_add;type(datetime)"`
}

func (m *User) TableName() string {
	return "user"
}

func init() {
	orm.RegisterModel(new(User))
}

//注册用户
func Register(regType string, user User) (int64, int64) {
	o := orm.NewOrm()
	var sql string
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, -1
	}
	user.Password = string(password)
	//payPassword, err := bcrypt.GenerateFromPassword([]byte(user.PayPassword), bcrypt.DefaultCost)
	//if err != nil {
	//	return 0, -2
	//}
	//user.PayPassword = string(payPassword)
	if regType == "email" {
		sql = "INSERT INTO `user`(`email`,`password`,`parent_id`,`parent_tree`,`country_id`,`mobile_code`,`nickname`,`level`) VALUES(?,?,?,?,?,?,?,?)"
		res, err := o.Raw(sql, user.Email, user.Password, user.ParentId, user.ParentTree, user.CountryId, user.MobileCode, user.Nickname, user.Level).Exec()
		if err == nil {
			id, _ := res.LastInsertId()
			return id, 0
		}
	} else if regType == "mobile" {
		sql = "INSERT INTO `user`(`mobile`,`password`,`parent_id`,`parent_tree`,`country_id`,`mobile_code`,`nickname`,`level`) VALUES(?,?,?,?,?,?,?,?)"
		res, err := o.Raw(sql, user.Mobile, user.Password, user.ParentId, user.ParentTree, user.CountryId, user.MobileCode, user.Nickname, user.Level).Exec()
		if err == nil {
			id, _ := res.LastInsertId()
			return id, 0
		}
	}

	return 0, -3
}

//获取邀请码
func GetInviteCode(id int64) string {
	o := orm.NewOrm()
	user := User{Id: id}
	err := o.Read(&user, "id")
	if err == nil {
		if strings.Trim(user.InviteCode, " ") != "" {
			return user.InviteCode
		}
		inviteCode := hash.IdToHash(int(id), "tet_invite", 6)
		result, _ := o.QueryTable("user").Filter("id", id).Update(orm.Params{
			"invite_code": inviteCode,
		})
		if result > 0 {
			return inviteCode
		}
	}
	return ""
}

//批量处理邀请码
func BatchDealInviteCode() bool {
	o := orm.NewOrm()
	sql := "SELECT id FROM `user` WHERE invite_code='' OR invite_code IS NULL;"
	var maps []orm.Params
	var id int64
	var inviteCode string
	o.Raw(sql).Values(&maps)
	if maps != nil && len(maps) > 0 {
		for _, v := range maps {
			id, _ = strconv.ParseInt(v["id"].(string), 10, 64)
			inviteCode = hash.IdToHash(int(id), "tet_invite", 6)
			sql = fmt.Sprintf("UPDATE `user` SET invite_code='%v' WHERE id=%v AND (invite_code='' OR invite_code IS NULL);", inviteCode, id)
			o.Raw(sql).Exec()
		}
		return true
	}
	return false
}

//通过邀请码获取用户ID
func GetUserByInviteCode(inviteCode string) (u User) {
	o := orm.NewOrm()

	inviteCode = strings.Trim(inviteCode, " ")
	if inviteCode == "" {
		return u
	}
	sql := "SELECT * FROM `user` WHERE invite_code=?;"
	o.Raw(sql, inviteCode).QueryRow(&u)
	return u
}

//通过用户名获取用户
func GetUserByUserName(userName string) User {
	o := orm.NewOrm()
	var user User
	_ = o.Raw("SELECT * FROM `user` WHERE `email` = ? OR `mobile`=?", userName, userName).QueryRow(&user)
	return user
}

//通过用户ID获取用户信息
func GetUserById(id int64) User {
	o := orm.NewOrm()
	var user User
	_ = o.Raw("SELECT * FROM `user` WHERE `id` = ?", id).QueryRow(&user)
	return user
}

//通过用户ID获取用户信息邀请人
func GetInviterById(id int64) orm.Params {
	o := orm.NewOrm()
	var user User
	var maps []orm.Params
	_ = o.Raw("SELECT * FROM `user` WHERE `id` = ?", id).QueryRow(&user)

	sql := "SELECT T0.`id`, T0.`avatar` , T0.`mobile` , T0.`mobile_code` , T0.`email` , T0.`country_id`   , T0.`parent_id` , T0.`parent_tree` , T0.`level` , T0.`invite_code` , T0.`status` , T0.`created_at` , T0.`updated_at`  FROM user T0 where T0.`id` =?;"
	o.Raw(sql, id).Values(&maps)
	if maps != nil {
		return maps[0]
	}

	return nil
}

//通过用户UnionID获取用户信息
func GetUserByUnionId(unionId string) User {
	o := orm.NewOrm()
	var user User
	_ = o.Raw("SELECT * FROM `user` WHERE `union_id` = ?", unionId).QueryRow(&user)
	return user
}

//是否存在手机号
func IsMobileExist(mobile string) bool {
	o := orm.NewOrm()
	exist := o.QueryTable("user").Filter("mobile", mobile).Exist()
	return exist
}

//是否存在邮箱
func IsEmailExist(email string) bool {
	o := orm.NewOrm()
	exist := o.QueryTable("user").Filter("email", email).Exist()
	return exist
}

//找回密码
func FindPassword(findType string, user User) (string, bool) {
	o := orm.NewOrm()
	var (
		sql  string
		exit bool
	)
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", false
	}
	newPassWord := string(password)
	o.Begin()
	if findType == "email" {
		o.QueryTable("user").Filter("email", user.Email).One(&user)
		exit = o.QueryTable("user_withdraw_limit").Filter("user_id", user.Id).Exist()
		sql = "UPDATE `user` SET `password`=? WHERE email=? AND email IS NOT NULL;"
		rs, err := o.Raw(sql, newPassWord, user.Email).Exec()
		if err != nil {
			o.Rollback()
			return "", false
		} else if err == nil {
			temp, _ := rs.RowsAffected()
			if temp <= 0 {
				o.Rollback()
				return "", false
			}
		}
	} else if findType == "mobile" {
		o.QueryTable("user").Filter("mobile", user.Mobile).One(&user)
		exit = o.QueryTable("user_withdraw_limit").Filter("user_id", user.Id).Exist()
		sql = "UPDATE `user` SET `password`=? WHERE mobile=? AND mobile IS NOT NULL;"
		rs, err := o.Raw(sql, newPassWord, user.Mobile).Exec()
		if err != nil {
			o.Rollback()
			return "", false
		} else if err == nil {
			temp, _ := rs.RowsAffected()
			if temp <= 0 {
				o.Rollback()
				return "", false
			}
		}
	}
	if exit {

		sql2 := "UPDATE user_withdraw_limit set limit_at=DATE_SUB(CURRENT_TIMESTAMP,INTERVAL -1 DAY) where user_id=? "
		rs, err := o.Raw(sql2, user.Id).Exec()
		if err != nil {
			o.Rollback()
			return "", false
		} else if err == nil {
			temp, _ := rs.RowsAffected()
			if temp <= 0 {
				o.Rollback()
				return "", false
			}
		}
	} else {
		sql1 := "INSERT INTO user_withdraw_limit(`user_id`,`limit_at`) VALUES (?,DATE_SUB(CURRENT_TIMESTAMP,INTERVAL -1 DAY))"
		rs, err := o.Raw(sql1, user.Id).Exec()
		if err != nil {
			o.Rollback()
			return "", false
		} else if err == nil {
			temp, _ := rs.RowsAffected()
			if temp <= 0 {
				o.Rollback()
				return "", false
			}
		}
	}
	o.Commit()
	return "", true
}

//修改密码
func ChangePassword(user User) bool {
	o := orm.NewOrm()
	var sql string
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return false
	}
	user.Password = string(password)

	sql = "UPDATE `user` SET `password`=? WHERE id=?;"
	_, err = o.Raw(sql, user.Password, user.Id).Exec()
	if err == nil {
		//num, _ := res.RowsAffected()
		return true
	}
	return false
}

//修改支付密码
func ChangePayPassword(user User) bool {
	o := orm.NewOrm()
	var sql string
	password, err := bcrypt.GenerateFromPassword([]byte(user.PayPassword), bcrypt.DefaultCost)
	if err != nil {
		return false
	}
	user.PayPassword = string(password)

	sql = "UPDATE `user` SET `pay_password`=? WHERE id=?;"
	_, err = o.Raw(sql, user.PayPassword, user.Id).Exec()
	if err == nil {
		//num, _ := res.RowsAffected()
		return true
	}
	return false
}

//用户基本信息
func GetUserInfo(userId int64) UserInfo {
	o := orm.NewOrm()
	var userInfo UserInfo
	sql := "SELECT a.`id`,a.`nickname`,`avatar`,`mobile`,a.`mobile_code`,`email`,a.`country_id`,`country`,`country_en`,`parent_id`,`level`,`invite_code`,a.`status`,a.`head_img`,IFNULL(c.`status`,-1) AS auth_status FROM `user` a " +
		" LEFT JOIN country b ON (a.country_id=b.id) " +
		" LEFT JOIN user_auth c ON (a.id=c.user_id) " +
		" WHERE a.id=?"
	err := o.Raw(sql, userId).QueryRow(&userInfo)
	fmt.Println(err)
	return userInfo
}

//绑定邮箱
func BindEmail(email string, user User) (bool, int64) {
	o := orm.NewOrm()
	o.Begin()
	var (
		sql string
		err error
	)
	sql = "UPDATE `user` SET `email`=? WHERE id=?;"
	_, err = o.Raw(sql, email, user.Id).Exec()
	if err != nil {
		o.Rollback()
		return false, -1
	}
	sql = "INSERT INTO `user_bind_change` (`user_id`,`type`,`old_value`,`content`)VALUES(?,'email',?,?);"
	content := fmt.Sprintf("用户：%v修改绑定邮箱为：%v，原邮箱为%v", user.Id, email, user.Email)
	_, err = o.Raw(sql, user.Id, user.Email, content).Exec()
	if err != nil {
		o.Rollback()
		return false, -2
	}

	o.Commit()
	return true, 0
}

//绑定手机修改
func BindMobile(mobile string, mobileCode string, countryId int64, user User) (bool, int64) {
	o := orm.NewOrm()
	o.Begin()
	var (
		sql string
		err error
	)
	sql = "UPDATE `user` SET `mobile`=?,`mobile_code`=?,`country_id`=? WHERE id=?;"
	_, err = o.Raw(sql, mobile, mobileCode, countryId, user.Id).Exec()
	if err != nil {
		o.Rollback()
		return false, -1
	}
	sql = "INSERT INTO `user_bind_change` (`user_id`,`type`,`old_value`,`content`)VALUES(?,'mobile',?,?);"
	oldValue := user.MobileCode + user.Mobile
	content := fmt.Sprintf("用户：%v修改绑定手机号为：%v，手机码为：%v，原手机号为%v", user.Id, mobile, mobileCode, oldValue)
	_, err = o.Raw(sql, user.Id, oldValue, content).Exec()
	if err != nil {
		o.Rollback()
		return false, -2
	}

	o.Commit()
	return true, 0
}

//修改用户昵称
func UpdateUserNickNameById(id int64, newNickName string) int64 {
	o := orm.NewOrm()
	result, _ := o.QueryTable("user").Filter("id", id).Update(orm.Params{
		"nickname": newNickName,
	})
	return result
}

//修改临时表的昵称
func UpdateUserTeamTempleByUserId(userId int64, newNickName string) int64 {
	o := orm.NewOrm()
	o.Raw("update team_temple set nickname=? where user_id=?", newNickName, userId).Exec()
	return 1
}

// 修改用户信息父级
func UpdateUserParentIdParentTreeAndLevel(userId, parentId int64) int64 {
	o := orm.NewOrm()
	result, _ := o.QueryTable("user").Filter("id", userId).Update(orm.Params{
		"parent_id":  parentId,
		"invited_at": time.Now(),
	})

	//更新用户树
	_, err := o.Raw("UPDATE `user` SET parent_tree=getParentTree(id)  where FIND_IN_SET(?,parent_tree) or id=?", userId, userId).Exec()
	if err != nil {
		//o.Rollback()
		return 0
	}
	//更新level
	sql1 := "UPDATE `user` SET `level`=1 WHERE parent_id>0 AND parent_id=parent_tree AND FIND_IN_SET(?,parent_tree);"
	o.Raw(sql1, parentId).Exec()
	sql1 = "UPDATE `user` SET `level`=LENGTH(parent_tree)-LENGTH(REPLACE(`parent_tree`,',',''))+1 WHERE parent_id>0 AND parent_id!=parent_tree AND FIND_IN_SET(?,parent_tree);"
	o.Raw(sql1, parentId).Exec()

	return result
}

//修改头像
func UpdateUserHeadImg(user User) (string, bool) {
	var (
		msg string
	)
	o := orm.NewOrm()
	o.Begin()
	i, err := o.Update(&user, "head_img")
	if err != nil || i <= 0 {
		msg = "request_fail"
		o.Rollback()
		return msg, false
	}
	msg = "request_success"
	o.Commit()
	return msg, true
}

//我的团队成员信息
func GetMyTeam(userId int64, page, pageSize int) tool.Pager {
	o := orm.NewOrm()
	var (
		list []orm.Params
		row  tool.Row
	)
	beginIndex := (page - 1) * pageSize
	result := tool.Pager{Page: page, PageSize: pageSize, TotalCount: row.RowCount}
	sql := "select T0.id,T0.nickname,T0.mobile,T0.email,T0.created_at  from user T0 where  T0.parent_id = ? ORDER BY updated_at DESC LIMIT ?,?"
	num, err := o.Raw(sql, userId, beginIndex, pageSize).Values(&list)
	if err == nil && num > 0 {
		result.List = list
		result.TotalCount = num
	}
	return result

}
