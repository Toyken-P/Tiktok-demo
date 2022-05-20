package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"path/filepath"
	"time"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	db, err := gorm.Open("mysql", "root:123456@tcp(121.43.190.40:3306)/douyin?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		fmt.Println("connect db failed err:", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "网络出现了点小问题，请稍后再试",
		})
		return
	}
	token := c.PostForm("token")
	var users = make([]User, 999)
	var select_user User
	flag := false
	db.Find(&users)
	for _, user := range users {
		if user.UserName+user.Password == token {
			flag = true
			select_user = user
		}
	}
	if !flag {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "please login first !",
		})
		return
	}
	data, err := c.FormFile("data")
	if err != nil {
		fmt.Println("get file failed err:", err)
		c.JSON(
			http.StatusOK,
			Response{
				StatusCode: 1,
				StatusMsg:  "upload failed ",
			})
		return
	}
	filename := filepath.Base(data.Filename)
	playerUri := playerprifex + filename
	saveFile := "./public/videos/" + filename
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	var count_video Video
	//db.Table("videos").Find(&count)
	db.Last(&count_video)
	createtime := time.Now().Unix()
	title := c.Request.PostForm.Get("title")
	video := Video{
		Id:        count_video.Id + 1,
		Author:    select_user,
		PlayUrl:   playerUri,
		CoverUrl:  "",
		CreatedAt: createtime,
		Title:     title,
	}
	db.Create(video)
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  filename + " upload successfully",
	})

}

func PublishList(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	if token == "" {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "please login first",
			},
			VideoList: nil,
		})
	}
	db, err := gorm.Open("mysql", "root:123456@tcp(121.43.190.40:3306)/douyin?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		fmt.Println("connect db failed err:", err)
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
			},
		})
		return
	}
	var users = make([]User, 999)
	db.Find(&users)
	var user_id int64
	for _, user := range users {
		if token == user.UserName+user.Password {
			user_id = user.Id
		}
	}
	if user_id == 0 {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
			},
		})
	}
	videolist := make([]Video, 999)
	db.Where("user_id = ?", user_id).Find(&videolist)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videolist,
	})

}
