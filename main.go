package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	conf              *config = &config{}
	configPath        *string = nil
	defaultConfigPath string  = "config.sample.yaml"
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
	parseConf()
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

func parseConf() {
	configPath = flag.String("c", "config.yaml", "path to config.yaml")
	flag.Parse()
	// config.yaml不存在则使用config.sample.yaml
	configFile, err := os.OpenFile(*configPath, os.O_RDWR, 0644)
	if err != nil {
		switch {
		case os.IsPermission(err):
			checkError(fmt.Errorf("无法访问配置文件: %s", *configPath))
		case os.IsNotExist(err):
			// 打开样板配置文件
			sampleConfigFile, err := os.OpenFile(defaultConfigPath, os.O_RDWR, 0644)
			checkError(err)
			defer sampleConfigFile.Close()
			// 创建配置文件
			configFile, err = os.OpenFile(*configPath, os.O_CREATE|os.O_RDWR, 0644)
			checkError(err)
			defer configFile.Close()
			// 复制: 样板配置 => 配置
			confByteNum, err := io.Copy(configFile, sampleConfigFile)
			checkError(err)
			if confByteNum < 1 {
				checkError(fmt.Errorf("复制%s到%s失败", defaultConfigPath, *configPath))
			}
			log.Printf("使用默认配置文件: %s (%d字节)", defaultConfigPath, confByteNum)
			// 写入磁盘
			err = configFile.Sync()
			checkError(err)
			// 从头读取
			_, err = configFile.Seek(0, io.SeekStart)
			checkError(err)
		default:
			checkError(err)
		}
	} else {
		defer configFile.Close()
	}
	// 读取配置文件
	configStr, err := ioutil.ReadAll(configFile)
	checkError(err)
	// 解析配置文件
	err = yaml.Unmarshal(configStr, conf)
	checkError(err)
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
