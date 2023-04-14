# httphelper

`httphelper` 是一个 Golang 包，提供了一些实用的 HTTP 请求辅助函数和抽象接口，用于简化 HTTP 请求的构建和处理过程。

## 安装

使用 `go get` 命令安装：


`go get github.com/artisancloud/httphelper`

## 功能

- 支持设置请求头、查询参数、请求体等
- 支持 multipart/form-data 类型的请求体构建
- 支持 JSON、XML 等类型的请求体构建
- 支持中间件，可以自定义处理请求前和请求后的逻辑
- 支持返回结果自动解析为指定的类型

## 使用示例

```
package main

import (
    "fmt"
    "github.com/artisancloud/httphelper"
    "github.com/artisancloud/httphelper/client"
)

func main() {
    conf := &httphelper.Config{
        Config: &client.Config{
            Timeout: 30,
        },
        BaseUrl: "http://localhost:8000",
    }

	helper, _ := httphelper.NewRequestHelper(conf)

	res, err := helper.Df().
		Method("GET").
		Uri("/api/users").
		Query("page", "1").
		Query("size", "10").
		Request()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res.Status)

}
```

## 贡献

欢迎对本项目提交问题和建议，也欢迎参与项目的开发。

_README.md by ChatGPT_