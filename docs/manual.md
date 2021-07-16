# 服务端开发手册

## 下载go及配置环境
1. mac下先安装包管理工具 [brew](https://brew.sh)
2. 通过brew安装go，需要1.16以上: 
   ``` shell
   brew install go
   ```
3. 配置go的环境
	* 建立一个go的工作目录，比如 `~/go/`，工作目录下建立三个子目录分别为
  		```
		  src // 克隆代码都在src目录内
		  pkg // go module获取的代码都在pkg目录内
		  bin // go相关工具都在bin目录内
		```
	* 配置环境，终端中运行以下命令：
  		``` shell
		  // 配置代理，否则有些包无法下载
		  export GOPROXY=https://goproxy.io,direct 

		  // 配置go的工作路径，使用`go get`命令时会将代码和工具安装到对应的路径
		  export GOPATH=$HOME/go

		  // 配置go工具的安装目录
		  export GOBIN=$GOPATH/bin

		  // 配置私有仓库路径
		  export GOPRIVATE=e.coding.net

		  // 配置shell的PATH
		  export PATH=$PATH:$GOPATH:$GOBIN
		```

---------
  
## 依赖中间件
* [mongodb](https://www.mongodb.com)
* [protobuf](https://developers.google.com/protocol-buffers)
* [go-micro@v3](https://github.com/asim/go-micro)

----------

## 软件安装
1. mac下安装mongodb官方教程 [https://docs.mongodb.com/manual/tutorial/install-mongodb-on-os-x/](https://docs.mongodb.com/manual/tutorial/install-mongodb-on-os-x/)
2. mac下安装protobuf: 
	``` shell
	brew install protobuf
	```
3. mac下安装protoc-gen-go: 
   ``` shell
   go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
   ```
4. mac下安装protoc-gen-micro: 
   ``` shell
   go get github.com/asim/go-micro/cmd/protoc-gen-micro/v3
   ```

----------------
## 服务开发流程
1. 克隆需要的代码仓库：
   ``` shell
   git clone https://github.com/east-eden/server.git $GOPATH/src/github.com/east-eden/server 

   git clone https://github.com/east-eden/proto.git $GOPATH/src/github.com/east-eden/proto 

   git clone https://github.com/east-eden/excel.git $GOPATH/src/github.com/east-eden/excel 

   git clone https://github.com/east-eden/server_bin.git $GOPATH/src/github.com/east-eden/server_bin 
   ```
2. 服务端开发所需仓库组织结构图：
   ```
	.
	├── excel 			// 策划表及导出工具
	├── proto 			// proto协议定义及生成代码工具
	├── server			// 服务端所有代码
	└── server_bin 		// 服务端打包版本
   ```

3. 服务端代码结构图：
   ```
   .
	├── Jenkinsfile 				// Jenkins pipeline
	├── LICENSE
	├── Makefile 					// make指令
	├── README.md
	├── apps 						// 所有服务入口及Dockerfile
	├── bitbucket-pipelines.yml		// bitbucket pipeline
	├── ci-building-base.Dockerfile // coding 基础镜像
	├── config						// 所有服务及中间件配置文件
	├── data						// 所有服务、中间件以及log持久化的文件
	├── define 						// 所有的常量定义
	├── docker-compose.yaml			// docker-compose配置文件
	├── docs						// 文档
	├── excel						// excel相关工具
	├── go.mod
	├── go.sum
	├── internal					// 内部包
	├── logger
	├── proto						// protobuf生成的代码
	├── services					// 所有服务代码
	├── store						// 存储相关
	├── transport					// 协议相关
	├── utils						// 常用工具
	└── vendor						// 依赖的模块代码
   ```

4. 开启服务
    * 开启网关服务：
	``` shell
	cd apps/gate
	go run main.go
	```

	* 开启游戏逻辑服务：
	``` shell
	cd apps/game
	go run main.go
	```

	* 开启邮件服务：
	``` shell
	cd apps/mail
	go run main.go
	```

	* 开启cli客户端:
	``` shell
	cd apps/client
	go run main.go
	```

-------------
## excel 导出工具使用
* svn或者git更新了策划表之后，在`$GOPATH/src/github.com/east-eden/excel/`路径下使用命令`make gen_mac`，即可根据最新的策划excel表格导出对应的数据结构到目录`$GOPATH/src/github.com/east-eden/server/excel/auto`中，并且导出对应的数据到目录`$GOPATH/src/github.com/east-eden/server/config/csv/`中，服务开启时会加载`csv/`路径下的数据文件

## proto 导出工具使用
* git更新了`*.proto`文件后，在`$GOPATH/src/github.com/east-eden/proto/`路径下使用命令`make proto`，即可生成对应的`*.pb.go`和`*.pb.micro.go`文件到路径`$GOPATH/src/github.com/east-eden/server/proto/`中