package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type smzdm interface {
	getPostID() []int
	smzdmSign() (bool, error)
	smzdmCommit(int, string) (bool, error)
}

type cracker struct {
	account account
}

// NewCracker ...
func NewCracker(account account) *cracker {
	return &cracker{account: account}
}

func (c *cracker) handle(method, url, referer string, body *io.Reader) []byte {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	checkError(err)
	header := http.Header{}
	if referer == "" {
		referer = "https://www.smzdm.com"
	}
	header.Add("Referer", referer)
	header.Add("User-Agent", c.account.UserAgent)
	header.Add("Cookie", c.account.Cookies)
	req.Header = header
	resp, err := client.Do(req)
	checkError(err)
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	return content
}
func (c *cracker) getPostID() []int {
	postID := make([]int, conf.PostCommentMax)
	mtrand := rand.New(rand.NewSource(time.Now().UnixNano()))
	randNum := int(30.0 + math.Round(mtrand.Float64()*40.0))
	url := "https://faxian.smzdm.com/h1s0t0f37c0p" + strconv.Itoa(randNum)
	content := c.handle("GET", url, "https://www.smzdm.com/jingxuan/", nil)
	reg, err := regexp.Compile(`articleid="\d_(\d+)"`)
	checkError(err)
	postIDMatch := reg.FindAllSubmatch(content, conf.PostCommentMax)
	for k, v := range postIDMatch {
		for k2, v2 := range v {
			if k2 == 1 {
				id, err := strconv.Atoi(string(v2))
				checkError(err)
				postID[k] = id
			}
		}
	}
	return postID
}

func (c *cracker) smzdmSign() (bool, error) {
	jsonData := &signJson{}
	ts := time.Now().UnixNano()
	url := fmt.Sprintf("https://zhiyou.smzdm.com/user/checkin/jsonp_checkin?callback=jQuery112409568846254764496_%d&_=%d", ts, ts)
	content := c.handle("GET", url, "http://www.smzdm.com/qiandao/", nil)
	reg := regexp.MustCompile(`^jQuery\d+_\d+\((.*?)\)$`)
	jsonStr := reg.ReplaceAll(content, []byte(`$1`))
	json.Unmarshal(jsonStr, jsonData)
	if jsonData.ErrorCode == 0 {
		signResult = append(signResult, *jsonData)
		return true, nil
	}
	return false, errors.New(string(jsonStr))
}

func (c *cracker) smzdmCommit(postID int, comment string) (bool, error) {
	jsonData := &commentJson{}
	ts := time.Now().UnixNano()
	url := fmt.Sprintf("https://zhiyou.smzdm.com/user/comment/ajax_set_comment?callback=jQuery111006551744323225079_%d&type=3&pid=%d&parentid=0&vote_id=0&vote_type=&vote_group=&content=%s&_=%d", ts, postID, comment, ts)
	content := c.handle("GET", url, "https://zhiyou.smzdm.com/user/submit/", nil)
	reg := regexp.MustCompile(`^jQuery\d+_\d+\((.*?)\)$`)
	jsonStr := reg.ReplaceAll(content, []byte(`$1`))
	json.Unmarshal(jsonStr, jsonData)
	if jsonData.ErrorCode == 0 {
		commentResult = append(commentResult, *jsonData)
		return true, nil
	}
	return false, errors.New(string(jsonStr))
}
