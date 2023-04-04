package schemas

type Swipe struct {
	Swiper int `json:"swiper"`
	Swipee int `json:"swipee"`
	Comment string `json:"comment"`
	Like bool `json:"like"`
}