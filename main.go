package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"./internal/tformatter"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type Stub struct {
	FlagName  string
	FlagValue string
	FlagUsage string
}

type Cmd struct {
	CmdName    string
	Parameters string
	CmdHandler Handler
}

var (
	cmds = []Cmd{
		{"bucket_list", "", ListBucket},
		{"bucket_new", "bucket_name", NewBucket},
		{"put", "bucket_name,key,file", PutObject},
		{"list", "bucket_name", ListObject},
		{"get", "bucket_name,key,file", GetObject},
		{"input", "bucket_name,key,c", InputObject},
		{"output", "bucket_name,key", OutputObject},
		{"delete", "bucket_name,key", DeleteObject},
	}
	stubs = []Stub{
		{"bucket_name", "lewiskong", "bucket name"},
		{"key", "", "object key"},
		{"file", "", "object value(file path)"},
		{"c", "", "object value(content)"},
	}
)

func DeleteObject() error {
	c := NewClient(*params["bucket_name"])
	name := *params["key"]
	if name == "" {
		return errors.New("key cannot be nil")
	}

	_, err := c.Object.Delete(context.Background(), name)
	return err
}

func OutputObject() error {
	c := NewClient(*params["bucket_name"])
	name := *params["key"]
	if name == "" {
		return errors.New("key cannot be nil")
	}

	resp, err := c.Object.Get(context.Background(), name, nil)
	if err != nil {
		return err
	}

	bts, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Print(string(bts[:]))
	return nil
}

func GetObject() error {
	c := NewClient(*params["bucket_name"])
	name, file := *params["key"], *params["file"]
	if name == "" {
		return errors.New("key cannot be nil")
	}
	if file == "" {
		return errors.New("file path cannot be nil")
	}

	_, err := c.Object.GetToFile(context.Background(), name, file, nil)
	return err
}

func ListObject() error {
	c := NewClient(*params["bucket_name"])
	opt := &cos.BucketGetOptions{
		//Prefix:  "test",
		MaxKeys: 3,
	}
	v, _, err := c.Bucket.Get(context.Background(), opt)
	if err != nil {
		panic(err)
	}

	fout := tfmt.New("name", "size", "last_modified", "url")
	for _, c := range v.Contents {
		dpath := fmt.Sprintf("https://%s-1253808810.cos.ap-guangzhou.myqcloud.com/%s", *params["bucket_name"], c.Key)
		t, _ := time.Parse("2006-01-02T15:04:05.000Z", c.LastModified)
		fout.Println(c.Key, c.Size, t.Local().Format("2006-01-02 15:04:05"), dpath)
	}
	fout.Output()

	return nil
}

func InputObject() error {
	c := NewClient(*params["bucket_name"])
	name, content := *params["key"], *params["c"]
	if name == "" {
		return errors.New("key cannot be nil")
	}
	if content == "" {
		return errors.New("content cannot be nil")
	}

	_, err := c.Object.Put(context.Background(), name, strings.NewReader(content), nil)
	return err
}

func PutObject() error {
	c := NewClient(*params["bucket_name"])
	name, fpath := *params["key"], *params["file"]
	if fpath == "" {
		return errors.New("file path cannot be nil")
	}
	if name == "" {
		if name = getFileName(fpath); name == "" {
			return errors.New("key cannot be nil")
		}
	}

	_, err := c.Object.PutFromFile(context.Background(), name, fpath, nil)
	if err != nil {
		return err
	}
	fmt.Printf("https://%s-1253808810.cos.ap-guangzhou.myqcloud.com/%s\n", *params["bucket_name"], name)
	return nil
}

func NewBucket() error {
	c := NewClient(*params["bucket_name"])
	_, err := c.Bucket.Put(context.Background(), nil)
	return err
}

func ListBucket() error {
	bucketLevelClient := NewClient("")
	s, _, err := bucketLevelClient.Service.Get(context.Background())
	if err != nil {
		return err
	}
	for _, b := range s.Buckets {
		fmt.Printf("Name: %s\tRegion: %s\tCreateDate: %s.\n", b.Name, b.Region, b.CreationDate)
	}
	return nil
}

func main() {
	flag.Parse()
	var err error
	for _, c := range cmds {
		if c.CmdName == *cmd {
			if err = c.CmdHandler(); err != nil {
				panic(err.Error())
			} else {
				return
			}
		}
	}
	panic("command not find")
}

func NewClient(bucketName string) *cos.Client {
	if bucketName == "" {
		return cos.NewClient(nil, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  conf.SecretID,
				SecretKey: conf.SecretKey,
			},
		})
	}
	u, _ := url.Parse(fmt.Sprintf("http://%s-1253808810.cos.ap-guangzhou.myqcloud.com", bucketName))
	b := &cos.BaseURL{BucketURL: u}
	// cos永久密钥
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.SecretID,
			SecretKey: conf.SecretKey,
		},
	})
	return client
}

func init() {
	// 加载配置
	bts, err := ioutil.ReadFile("/Users/nirvana/common/conf/conf.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	if err = json.Unmarshal(bts, &conf); err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	// 初始化命令行参数
	var help string
	for _, c := range cmds {
		help += fmt.Sprintf("\t%s(%s)\n", c.CmdName, c.Parameters)
	}
	cmd = flag.String("cmd", "", help)
	for _, stub := range stubs {
		params[stub.FlagName] = flag.String(stub.FlagName, stub.FlagValue, stub.FlagUsage)
	}
}

var (
	params = make(map[string]*string)
	cmd    *string
	conf   struct {
		SecretID  string `json:"SecretId"`
		SecretKey string `json:"SecretKey"`
	}
)

type Handler func() error

func getFileName(name string) string {
	i := strings.LastIndex(name, "/")
	if i == -1 {
		return name
	}
	return name[i+1:]
}
