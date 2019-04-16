package baseClass

import (
	"time"

	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"

	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type User struct {
	Id       int
	UserName string `orm:"unique"` //用户名
	Password string //密码
	Power    int    //权限码
	Email    string //邮箱
}

type LoginRecord struct {
	Id        int
	LoginTime int64  //登陆时间
	LoginIP   string //登陆IP
	Password  string `orm:"null"`    //使用的错误的密码(正确时为空)
	User      *User  `orm:"rel(fk)"` //User表外键
}

func (c *User) Login_matching() int { //明文形式密码比对
	o := orm.NewOrm()
	var user User
	if err := o.QueryTable("user").Filter("user_name__exact", c.UserName).One(&user); err != nil {
		if err == orm.ErrNoRows {
			return 1
		} else {
			return 3
			beego.Error("username error : ", err)
		}
	}
	if c.UserName != user.UserName {
		return 1
	}
	if a := strings.Compare(Encrypt(c.Password), user.Password); a != 0 {
		return 2
	}
	return 0
}

func (c *User) Login_matching_crypte() bool { //session密码比对
	o := orm.NewOrm()
	var user User
	if err := o.QueryTable("user").Filter("user_name", c.UserName).One(&user); err != nil {
		beego.Error("username error : ", err)
		return false
	}
	if a := strings.Compare(c.Password, user.Password); a != 0 {
		return false
	}
	return true
}

func (c *User) Login_record(password, loginIp string) bool {
	var record LoginRecord
	record.LoginTime = time.Now().Unix()
	record.LoginIP = loginIp
	record.Password = password
	record.User = c
	if _, err := orm.NewOrm().Insert(&record); err != nil {
		return false
	}
	return true
}

//加密
//md5+sha256
func Encrypt(password string) (cryptograph string) {
	salt1 := "fi22.ij5.,2432!i"
	salt2 := "fo2.43o5h2f(juaz"
	md5_obj := md5.New()
	md5_obj.Write([]byte(salt1 + password + salt2))
	md5_encode := md5_obj.Sum(nil)
	md5_string := hex.EncodeToString(md5_encode)

	salt3 := "easfcvadwa"
	salt4 := "ofkafjdisa"
	sha256_obj := sha256.New()
	sha256_obj.Write([]byte(salt3 + md5_string + salt4))
	sha1_encode := sha256_obj.Sum(nil)
	cryptograph = hex.EncodeToString(sha1_encode)

	return cryptograph
}
