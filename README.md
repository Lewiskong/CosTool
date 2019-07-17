# CosTool
一个基于腾讯云COS的命令行工具

## Compile
  ```shell
    go build -o cos
  ```
## Usage
  ```
    cos
       -cmd string
            bucket_list()
            bucket_new(bucket_name)
            put(bucket_name,key,file)
            list(bucket_name)
            get(bucket_name,key,file)
            input(bucket_name,key,c)
            output(bucket_name,key)
            delete(bucket_name,key)
      -bucket_name string
            bucket name (default "lewiskong")
      -c string
            object value(content)
      -file string
            object value(file path)
      -key string
            object key

  ```
  
## 个人用于
  * 图片图床(CDN)
  * 常用字符/密码存储（macos可结合pbcopy）
  * 私人文件云存储
  
  
