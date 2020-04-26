package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	conf = &config{}
)

type account struct {
	Username  string `yaml:"username"`
	Phone     string `yaml:"phone"`
	UserAgent string `yaml:"userAgent"`
	Cookies   string `yaml:"cookies"`
}
type config struct {
	EmailName        string    `yaml:"emailName"`
	EmailPassword    string    `yaml:"emailPassword"`
	ToEmail          string    `yaml:"toEmail"`
	DefaultUserAgent string    `yaml:"defaultUserAgent"`
	DelayMin         float32   `yaml:"delayMin"`
	DelayMax         float32   `yaml:"delayMax"`
	Accounts         []account `yaml:"accounts"`
	Comments         []string  `yaml:"comments"`
}

func init() {
	configFile, e := os.OpenFile("config.yaml", os.O_RDWR, 0644)
	checkError(e)
	configStr, e := ioutil.ReadAll(configFile)
	checkError(e)
	e = yaml.Unmarshal(configStr, conf)
	checkError(e)
}
func main() {
	if len(conf.Accounts) > 0 {
		wg := &sync.WaitGroup{}
		for _, acnt := range conf.Accounts {
			wg.Add(1)
			acnt := acnt
			go func() {
				defer wg.Done()
				toSign(acnt)
			}()
		}
		wg.Wait()
	}
}

func toSign(acnt account) {
	// 实例化
	var smzdm smzdm = NewCracker(conf.Accounts[0])

	// ①签到
	ok, err := smzdm.smzdmSign()
	if !ok {
		panic(err)
	} else {
		log.Printf("签到成功")
	}

	// ②获取最新文章
	postID := smzdm.getPostID()
	log.Printf("获取最新文章: %+v\n", postID)

	mtrand := rand.New(rand.NewSource(time.Now().UnixNano()))

	// ③前5篇文章留言
	if len(postID) < 5 {
		log.Panicf("获取最新文章[]篇, 数据异常, 不予评论")
	}
	confCommentsLen := len(conf.Comments)
	for _, id := range postID {
		// 随机评论
		randCommentsID := int(mtrand.Float32() * float32(confCommentsLen-1))
		comment := conf.Comments[randCommentsID]
		// 防止发送过快
		td := time.Duration(conf.DelayMin+mtrand.Float32()*(conf.DelayMax-conf.DelayMin)) * time.Second
		log.Printf("文章[%d]将于[%d秒后]评论[%s]\n", id, td/time.Second, comment)
		time.Sleep(td)
		// 开始评论
		ok, err := smzdm.smzdmCommit(id, comment)
		if !ok {
			panic(err)
		} else {
			log.Printf("文章[%d]评论成功\n", id)
		}
	}
}
