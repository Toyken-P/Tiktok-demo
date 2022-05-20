package controller

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id     int64 `json:"id,omitempty"`
	Author User  `json:"author" gorm:"embedded"`
	//UserId        int64  `json:"user_id,omitempty"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	CreatedAt     int64  `json:"created_at,omitempty"`
	Title         string `json:"title,omitempty" gorm:"column:title"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user" gorm:"embedded"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id            int64  `json:"user_id,omitempty" gorm:"column:user_id"`
	Name          string `json:"name,omitempty"`
	UserName      string `json:"user_name,omitempty"`
	Password      string `json:"password,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}
type Loginer struct {
	Token  string `json:"token,omitempty"`
	UserId int64  `json:"userId,omitempty" gorm:"column:user_id"`
}

var (
	playerprifex = "http://localhost:8080/static/"
)
