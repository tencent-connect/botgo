
<h2>文件说明</h2>
```
.
├── conf //配置文件
│   ├── robot_remote.yaml
│   └── robot.yaml
├── go.mod
├── go.sum
├── handler.go //机器人各种事件处理入口
├── local //单实例模式，和机器人sdk对接的通用逻辑，一般不需要修改
│   └── run_robot.go
├── main.go
├── multi_robot
├── processor //实际业务逻辑处理器，业务改动主要是改动这层
│   ├── mock.go
│   └── processor.go
├── README.md
└── remote //多实例模式，和机器人sdk对接的通用逻辑，一般不需要修改
    └── run_robot.go

```
具体实现时，只需要改动main中的handler和processor中的处理逻辑即可

<h2>运行步骤</h2>
<h3>配置机器人信息</h3>
<p>在conf/robot.yaml里配置需要启动的机器人。</p>
<h3>编译</h3>
<code>go build</code>
<h3>启动</h3>
<code>./multi_robot</code>

<h3>单实例模式</h3>
```go
import "multi_robot/local"

func main(){
    ....

    local.InitProcessRobot("./conf/robot.yaml",
		processor.RegisterProcessor, processor.NewProcessor, []interface{}{
			InteractionEventHandler(),
			C2CMessageEventHandler(),
			GroupATMessageEventHandler(),
		}...)
    
    ....
}

```

<h3>分布式模式</h3>
```go
import "multi_robot/remote"

func main(){
    ....

    remote.InitProcessRobot("./conf/robot.yaml",
		processor.RegisterProcessor, processor.NewProcessor, []interface{}{
			InteractionEventHandler(),
			C2CMessageEventHandler(),
			GroupATMessageEventHandler(),
		}...)
    
    ....
}

```
配置文件里加上redis相关配置
```yaml
- appid: xxx
  appkey: xxxx
  redis:
     addr: xxxxxxxx
     net_work: 
     pass_word: 
  cluster_key: xxxxxxxxxx
- appid: xxxx
  appkey: xxxx
  redis:
     addr: xxxxxxxx
     net_work: 
     pass_word: 
  cluster_key: xxxxxxxxxxx

```