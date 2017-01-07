# Doob
![-_-](https://travis-ci.org/fudali113/Doob.svg?branch=master)

Doob is a rest and a simple router handler

init invoke AddHandlerFunc(forwardUrl,methodStr,func)
such as

* add static folder
```
Doob.AddStaicPrefix("/static")
```

* add forwardUrl func , use `&&` split urls
```
router := Doob.DefaultRouter()

router.AddHandlerFunc("/Doob/origin/{who}/{do} && /Doob/origin1/{who}/{do}", origin, Doob.GET, Doob.POST, Doob.PUT, Doob.DELETE)
```

* use `Get,Post,Put,Delete,Options,Head` method
```
  // get `/Doob` at the beginning base router
  doobPrefixRouter := Doob.GetRouter("Doob")

  // use `{}` distinction pathVariable
  doobPrefixRouter.Get("/{name}/{value}", func)
  // {} 中支持添加正则表达式，用`:`分割参数名和正则表达式
  doobPrefixRouter.Post("/{name:[0-9]+}/{value}", func)
  doobPrefixRouter.Put("/{name}/{value}", func)
  doobPrefixRouter.Delete("/{name}/{value}", func)
  doobPrefixRouter.Options("/{name}/{value}", func)
  doobPrefixRouter.Head("/{name}/{value}", func)
```

* support func classify
```
  // 兼容原始http方法类
  func origin(w http.ResponseWriter, r *http.Request) {}

  // 根据doob 里的context 进行获取参数或者返回
  func ctx(ctx *Doob.Context) interface{} {}

  // 根据url参数自动注入参数
  // 返回值为string时为返回静态文件
  // 返回不为string时默认将解析该对象，并返回给请求用户
  func di(name, value string) interface{} {}

  // 返回 string 是 type， interface{} 是你需要处理的数据
  // 你的返回值将流向注册的返回值处理器进行匹配并处理
  // 当只返回string时，把interface当做nil
  //
  // 当只返回数据时，此时returnDealDefaultType 为 auto
  // 此时将根据 請求的header Accept参数进行判断type类型
  // 你可以通过SetReturnDealDefaultType(t string)设置默认的type类型
  func returnHtml() (string, interface{}) {}
```

* next use
```
  Doob.Start(8888)
```

run you application

clone this project , run demo
```
  cd /sample
  go run demo.go
```
