# PGO2
PGO2应用框架即"Pinguo GO application framework 2.0"，是Camera360服务端团队基于PGO研发的一款简单、高性能、组件化的GO应用框架。受益于GO语言高性能与原生协程，使用PGO2可快速地开发出高性能的web应用程序。

参考文档：[pgo2-docs](https://github.com/pinguo/pgo2-docs)

应用示例：[pgo2-demo](https://github.com/pinguo/pgo2-demo)

## 基准测试

## 环境要求
- GO 1.13+
- Make 3.8+
- Linux/MacOS/Cygwin
- GoLand 2019 (建议)

## 项目目录
规范：
- 基于go mod 按go规范执行。

```
<project>
├── bin/                # 编译程序目录
├── configs/            # 配置文件目录
│   ├── production/     # 环境配置目录
│   │   ├── app.yaml
│   │   └── params.yaml
│   ├── testing/
│   ├── app.yaml        # 项目配置文件
│   └── params.yaml     # 自定义配置文件
├── makefile            # 编译打包
├── runtime/            # 运行时目录 主要用于存储日志
├── assets/             # 非web相关资源
├── web/                # 静态资源目录
│   ├── static          # 静态资源
│   ├── template        # 视图模板目录
├── cmd/                # 项目入口目录
└── pkg/                # 项目源码目录
    ├── command/        # 命令行控制器目录
    ├── controller/     # HTTP控制器目录
    ├── lib/            # 项目基础库目录
    ├── model/          # 模型目录(数据交互)
    ├── service/        # 服务目录(业务逻辑)
    ├── struct/         # 结构目录(数据定义)

```

## 快速开始
1. 拷贝makefile
    
    非IDE环境(命令行)下，推荐使用make做为编译打包的控制工具，从[pgo2](https://github.com/pinguo/pgo2)或[pgo2-demo](https://github.com/pinguo/pgo2-demo)将makefile复制到项目目录下。
    ```sh
    make start      # 编译并运行当前工程
    make stop       # 停止当前工程的进程
    make build      # 仅编译当前工程
    make update     # go mod get
    make install    # go mod download
    make pgo2       # 安装pgo2框架到当前工程
    make init       # 初始化工程目录
    make help       # 输出帮助信息
    ```

2. 创建项目目录(以下三种方法均可)
    - 执行`make init`创建目录
    - 参见《项目目录》手动创建
    - 从[pgo2-demo](https://github.com/pinguo/pgo2-demo)克隆目录结构

3. 修改配置文件(conf/app.yaml)
    ```yaml
    name: "pgo-demo"
    GOMAXPROCS: 2
    runtimePath: "@app/runtime"
    publicPath: "@app/web/static"
    viewPath: "@app/web/template"
    server:
        httpAddr: "0.0.0.0:8000"
        readTimeout: "30s"
        writeTimeout: "30s"
    components:
        log:
            levels: "ALL"
            targets:
                info:
                    levels: "DEBUG,INFO,NOTICE"
                    filePath: "@runtime/info.log"
                error:
                    levels: "WARN,ERROR,FATAL"
                    filePath: "@runtime/error.log"
                console: 
                    levels: "ALL"
    ```

4. 安装PGO(以下两种方法均可)
    - 在项目根目录执行`make pgo2`安装PGO2
    - 在项目根目录执行`go get -u github.com/pinguo/pgo`
5. 创建service(pkg/service/Welcome.go)
    ```go
    package Service

    import (
        "fmt"

        "github.com/pinguo/pgo2"
    )

    func NewWelcome() *Welcome{
       return &Welcome{}
    }

    type Welcome struct {
        pgo2.Object
    }
    
    func (w *Welcome) SayHello(name string, age int, sex string) {
        fmt.Printf("call in  service/Welcome.SayHello, name:%s age:%d sex:%s\n", name, age, sex)
    }
   
    ```
7. 创建控制器(pkg/controller/welcomeController.go)
    ```go
    package controller

    import (
        "service"
        "net/http"
     
        "github.com/pinguo/pgo"
    )

    type WelcomeController struct {
        pgo2.Controller
    }

    // 默认动作为index, 通过/welcome或/welcome/index调用
    func (w *WelcomeController) ActionIndex() {
        w.OutputJson("hello world", http.StatusOK)
    }
    
    // URL路由动作，根据url自动映射控制器及方法，不需要配置.
    // url的最后一段为动作名称，不存在则为index,
    // url的其余部分为控制器名称，不存在则为index,
    // 例如：/welcome/say-hello，控制器类名为
    // controller/welcomeController 动作方法名为ActionSayHello
    func (w *WelcomeController) ActionSayHello() {
        ctx := w.Context() // 获取PGO请求上下文件
    
        // 验证参数，提供参数名和默认值，当不提供默认值时，表明该参数为必选参数。
        // 详细验证方法参见Validate.go
        name := ctx.ValidateParam("name").Min(5).Max(50).Do() // 验证GET/POST参数(string)，为空或验证失败时panic
        age := ctx.ValidateQuery("age", 20).Int().Min(1).Max(100).Do() // 只验证GET参数(int)，为空或失败时返回20
        ip := ctx.ValidatePost("ip", "").IPv4().Do() // 只验证POST参数(string), 为空或失败时返回空字符串
    
        // 打印日志
        ctx.Info("request from welcome, name:%s, age:%d, ip:%s", name, age, ip)
        ctx.PushLog("clientIp", ctx.GetClientIp()) // 生成clientIp=xxxxx在pushlog中
    
        // 调用业务逻辑，一个请求生命周期内的对象都要通过GetObject()获取，
        // 这样可自动查找注册的类，并注入请求上下文(Context)到对象中。
        svc := w.GetObj(service.NewWelcome()).(*Service.Welcome)
    
        // 添加耗时到profile日志中
        ctx.ProfileStart("Welcome.SayHello")
        svc.SayHello(name, age, ip)
        ctx.ProfileStop("Welcome.SayHello")
    
        data := pgo.Map{
            "name": name,
            "age": age,
            "ip": ip,
        }
    
        // 输出json数据
        w.OutputJson(data, http.StatusOK)
    }
    
    // 正则路由动作，需要配置Router组件(components.router.rules)
    // 规则中捕获的参数通过动作函数参数传递，没有则为空字符串.
    // eg. "^/reg/eg/(\\w+)/(\\w+)$ => /welcome/regexp-example"
    func (w *WelcomeController) ActionRegexpExample(p1, p2 string) {
        data := pgo.Map{"p1": p1, "p2": p2}
        w.OutputJson(data, http.StatusOK)
    }
    
    // RESTFULL动作，url中没有指定动作名，使用请求方法作为动作的名称(需要大写)
    // 例如：GET方法请求ActionGET(), POST方法请求ActionPOST()
    func (w *WelcomeController) ActionGET() {
        w.GetContext().End(http.StatusOK, []byte("call restfull GET"))
    }
    
    ```
9. 创建程序入口(src/Main/main.go)
    ```go
    package main

    import (
        _ "controller" // 导入控制器

        "github.com/pinguo/pgo2"
    )

    func main() {
        pgo.Run() // 运行程序
    }
    ```
10. 编译运行
    ```sh
    make update
    make start
    curl http://127.0.0.1:8000/welcome
    ```

### 其它
参见[pgo2-docs](https://github.com/pinguo/pgo2-docs)
