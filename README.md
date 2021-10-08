# livego rtmp 流媒体服务

步态系统实时主要预览，目前只支持 H264

## 框架

> 主要包括 configure | av | container | parser | srv

### 配置 configure

> 服务配置初始化  
> 先加载默认配置，再读取配置文件，合并配置。

### 音视频基础 av

主要包括音视频处理基础接口定义

### 音视频打包容器 container

> flv：flv 格式容器编解码  
> ts：ts 格式容器编码

### 音视频编解码 parser

> aac： aac 音频格式编解码  
> mp3：mp3 音频格式编解码
> h264：h264 视频格式编解码

### 网络协议 protocol

amf | api | hls | http-flv | rtmp

> api： 实现系统 http 接口  
> 其他类似

### 服务接口 srv

> 实现 servicelib 接口，用于服务注册初始化配置中心配置等工作。  
> srv/version.go 实现版本升级功能

### 工具类 utils

livego 用到的基础工具类包
