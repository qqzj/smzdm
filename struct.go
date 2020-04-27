package main

type account struct {
	Username  string `yaml:"username"`
	Phone     string `yaml:"phone"`
	UserAgent string `yaml:"userAgent"`
	Cookies   string `yaml:"cookies"`
}
type config struct {
	EmailFrom         string    `yaml:"emailFrom"`
	EmailFromPassword string    `yaml:"emailFromPassword"`
	EmailFromSMTP     string    `yaml:"emailFromSMTP"`
	EmailToSubject    string    `yaml:"emailToSubject"`
	EmailTo           []string  `yaml:"emailTo"`
	DefaultUserAgent  string    `yaml:"defaultUserAgent"`
	DelayMin          float32   `yaml:"delayMin"`
	DelayMax          float32   `yaml:"delayMax"`
	PostCommentMax    int       `yaml:"postCommentMax"`
	Accounts          []account `yaml:"accounts"`
	Comments          []string  `yaml:"comments"`
}

type commonData struct {
	ErrorCode int `json:"error_code"`
	ErrorMsg  int `json:"error_msg"`
}

type signData struct {
	AddPoint    int    `json:"add_point"`
	CheckinNum  int    `json:"checkin_num"`
	Point       int    `json:"point"`
	Exp         int    `json:"exp"`
	Gold        int    `json:"gold"`
	Prestige    int    `json:"prestige"`
	Rank        int    `json:"rank"`
	Slogan      string `json:"slogan"`
	Cards       int    `json:"cards"`
	CanContract int    `json:"can_contract"`
}

type signJson struct {
	ErrorCode int      `json:"error_code"`
	ErrorMsg  int      `json:"error_msg"`
	Data      signData `json:"data"`
}

type commentData struct {
	PostPoints       int           `json:"post_points"`
	PostExperience   int           `json:"post_experience"`
	PostGold         int           `json:"post_gold"`
	PostPrestige     int           `json:"post_prestige"`
	CommentID        int           `json:"comment_ID"`
	FormatDate       string        `json:"format_date"`
	FormatDateClient string        `json:"format_date_client"`
	ParentData       []interface{} `json:"parent_data"`
	SinaUri          string        `json:"sina_uri"`
	CommentContent   string        `json:"comment_content"`
	IsAnonymous      int           `json:"is_anonymous"`
	DisplayName      string        `json:"display_name"`
	Head             string        `json:"head"`
}

type commentJson struct {
	ErrorCode int      `json:"error_code"`
	ErrorMsg  signData `json:"error_msg"`
}
