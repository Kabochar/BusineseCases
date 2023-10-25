# 项目架构

## 目录结构
- [ ] [main.go](#main)                // 系统启动入口
- [ ] [init](#init)                   // 系统启动初始化数据
- [ ] [auth](#auth)                   // 系统安全
- [ ] [file](#file)                   // 文件管理模块
- [ ] [framework](#framework)         // 系统框架
- [ ] [controller](#controller)       // 控制器目录、接口控制
- [ ] [service](#service)             // 服务目录、接口对应逻辑代码
- [ ] [model](#model)                 // 模型目录、所有结构体- 
- [ ] [monitor](#monitor)             // 系统监控
- [ ] [util](#util)                   // 工具
- [ ] [3rd](#3rd)                     // 第三方集成
## init
### 系统启动初始化数据

## auth
### 系统安全

```text
1.RESTFul API与API版本控制
2.中间件与jwt实现统一鉴权
```
- filter目录： 包含过滤器函数，用于对请求进行预处理或后处理，例如验证Token等。
- middleware目录： 包含中间件函数，用于处理请求之间的逻辑，例如权限验证、日志记录等。
  ```text
  1、open  开放接口   
  2、token token验证
  3、pp    公私钥验证接口
  限制IP 、 访问次数 、 浏览器指纹技术、系统所有API入库
  ```
- token目录： 包含Token管理的相关代码，用于生成、验证和管理Token。
- jwt目录： 包含JWT验证的相关代码，用于验证请求中的JWT并提取用户信息。

## framework
### 系统框架

## controller
### 接口控制代码

## service
### 接口对应逻辑代码

## model
### 所有结构体定义

## file
### 文件管理模块
```text
 1、定时删除临时文件
 2、文件类型控制
 3、上传文件按照年月日文件夹保存
 4、获取文件通过 redis -> db
 5、可以通过文件类动态生成同的接口
```


## 3rd
### 第三方集成
- redis目录： redis 操作。
```text
  业务代码不允许直接操作redis、 必须通过 工具函数操作 确保访问链路的完整性
```
## monitor
### 系统监控代码




