package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {

	latest_time := c.PostForm("latest_time")
	if latest_time == "" {
		latest_time = strconv.FormatInt(time.Now().Unix(), 10)
	}
	db, err := gorm.Open("mysql", "root:123456@tcp(121.43.190.40:3306)/douyin?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	if err != nil {
		fmt.Println("connect db failed : err:", err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "网络出了一点小问题，请稍后再试哦"},
		})
		return
	}
	readvideos := make([]Video, 30)
	fmt.Println("last_time : <--------------------------", latest_time)
	db.Where("created_at < ?", latest_time).Order("created_at desc").Find(&readvideos)
	fmt.Println("videos", readvideos)
	next_time := readvideos[len(readvideos)-1].CreatedAt
	fmt.Println("nexttime:", next_time)
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: readvideos,
		NextTime:  next_time,
	})
}
