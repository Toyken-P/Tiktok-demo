package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	db, err := gorm.Open("mysql", "root:123456@tcp(121.43.190.40:3306)/douyin?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		fmt.Println("connect db failed err:", err)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "网络出了一点小问题，请稍后再试哦"},
		})
		return
	}
	var user = User{UserName: username, Password: password}
	var aluser User
	db.Where(user).Find(&aluser)
	if aluser.Id != 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		return
	}
	if len(username) == 0 || len(password) == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "username or password can't be null !"},
		})
		return
	}
	var count_user User
	//db.Table("users").Count(&count)
	db.Last(&count_user)
	user.Id = count_user.Id + 1
	db.Create(&user)
	token := username + password
	db.Create(&Loginer{Token: token, UserId: user.Id})
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   user.Id,
		Token:    token,
	})

}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	var user = User{UserName: username, Password: password}
	db, err := gorm.Open("mysql", "root:123456@tcp(121.43.190.40:3306)/douyin?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		fmt.Println("connect db failed err:", err)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "网络出了一点小问题，请稍后再试哦"},
		})
		return
	}
	var select_user User
	db.Where(user).Find(&select_user)
	fmt.Println("select_user :", select_user)
	if select_user.Id == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户名或密码输入错误"},
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   select_user.Id,
			Token:    username + password,
		})
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	db, err := gorm.Open("mysql", "root:123456@tcp(121.43.190.40:3306)/douyin?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		fmt.Println("connect db failed err:", err)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "网络出了一点小问题，请稍后再试哦"},
		})
		return
	}
	var loginer Loginer = Loginer{Token: token}
	db.Where(loginer).Find(&loginer)
	if loginer.UserId != 0 {
		var user User
		db.Where("user_id = ?", loginer.UserId).Find(&user)
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
		return
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	})

}
