//go:generate goversioninfo -64 -o icon.syso
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
	startAt           time.Time = time.Now()
	conf              *config   = &config{}
	configPath        *string   = nil
	defaultConfigPath string    = "config.sample.yaml"
	signResult        []signJson
	commentResult     []commentJson
)

func init() {
	parseConf()
}

func main() {
	if len(conf.Accounts) > 0 {
		wg := &sync.WaitGroup{}
		for index, account := range conf.Accounts {
			wg.Add(1)
			// 实例化
			i := index
			a := account
			go func() {
				defer wg.Done()
				var smzdmer smzdm = NewCracker(i, a)
				toSignAndComment(smzdmer)
			}()
		}
		wg.Wait()
		send()
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

func toSignAndComment(smzdmer smzdm) {
	// ①签到
	ok, err := smzdmer.smzdmSign()
	if !ok {
		panic(err)
	} else {
		log.Printf("签到成功")
	}

	// ②获取最新文章
	postID := smzdmer.getPostID()
	log.Printf("获取最新文章: %+v\n", postID)

	// ③前n篇文章留言
	var usedCommentIDs map[int]bool
	usedCommentIDs = make(map[int]bool)
	mtrand := rand.New(rand.NewSource(time.Now().UnixNano()))
	if postIDNum := len(postID); postIDNum < conf.PostCommentMax {
		log.Panicf("获取最新文章[%d]篇, 数据异常, 不予评论", postIDNum)
	}
	confCommentsLen := len(conf.Comments)
	for _, id := range postID {
		// 随机评论
		randCommentsID := int(mtrand.Float32() * float32(confCommentsLen-1))
		// 检测这条评论是否已经用过了
		if _, ok := usedCommentIDs[randCommentsID]; ok {
			if randCommentsID < confCommentsLen-2 {
				randCommentsID++
			} else {
				randCommentsID--
			}
		}
		comment := conf.Comments[randCommentsID]
		// 防止发送过快
		td := time.Duration(conf.DelayMin+mtrand.Float32()*(conf.DelayMax-conf.DelayMin)) * time.Second
		log.Printf("文章[%d]将于[%d秒后]评论[%s]\n", id, td/time.Second, comment)
		time.Sleep(td)
		// 开始评论
		ok, err := smzdmer.smzdmCommit(id, comment)
		if !ok {
			panic(err)
		} else {
			log.Printf("文章[%d]评论成功\n", id)
		}
	}
}
