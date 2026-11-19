# 调试客户端与 GM 命令

`apps/client` 是仓库自带的调试客户端，启动后是**菜单式**交互（方向键/回车选择），
不是手敲命令行。下面按菜单说明它当前能做什么，以及 GM 命令的格式。

## 客户端菜单

主菜单：

| 菜单 | 说明 |
| --- | --- |
| 服务器连接管理 | 登录、心跳、断线重连 |
| 角色管理 | 创建角色、扫荡关卡、查询排行榜、执行 GM 命令 |
| 英雄管理 | 查看英雄、删除英雄 |
| 物品管理 | 查看物品、删除物品、使用物品 |
| 英雄碎片 | 用碎片合成英雄、用碎片合成收集品 |

### 服务器连接管理

| 菜单 | 说明 |
| --- | --- |
| 登录 | 输入 user id，走 TCP 登录 |
| websocket登录 | 输入 `user id,名字`，走 WebSocket 登录 |
| 发送心跳 | 发送 `C2S_HeartBeat` |
| 发送等待服务器返回消息 | 请求服务器时间等 |
| 客户端断开连接 | 客户端主动断开 |
| 服务器断开连接 | 请求服务端踢下线 |
| 重发登陆请求 | 重连后重发登录 |

### 角色管理

| 菜单 | 说明 |
| --- | --- |
| 创建角色 | 输入角色名创建角色 |
| 扫荡关卡 | 输入 `关卡id, 次数` |
| 请求排行榜 | 输入排行榜 id |
| gm命令 | 输入一条 GM 命令，格式见下 |

### 英雄 / 物品 / 碎片

| 菜单 | 说明 |
| --- | --- |
| 删除英雄 | 输入英雄配置 id |
| 删除物品 | 输入物品 id（物品以配置 id 为键） |
| 使用物品 | 输入物品 id |
| 英雄碎片合成 | 输入英雄配置 id，用 chip 合成英雄 |
| 收集品碎片合成 | 输入收集品 id |

## GM 命令

在「角色管理 → gm命令」里输入，第一个词是分组，后面是子命令：

| 分组 | 子命令 | 示例 |
| --- | --- | --- |
| `player` | `level`、`exp`、`vip`、`panic` | `gm player level 11` |
| `hero` | `add`、`exp`、`level`、`promote`、`frag` | `gm hero add 2090101`、`gm hero frag 2090101 100` |
| `item` | `add`、`delete` | `gm item add 10100000 1000` |
| `stage` | `pass`（同 `challenge`/`win`）、`sweep` | `gm stage pass 1`、`gm stage sweep 1 5` |
| `tower` | `floor`（同 `pass`）、`settle` | `gm tower floor 10` |
| `quest` | `reward` | `gm quest reward` |
| `collection` | `frag`、`active` | `gm collection frag 1 10` |
| `mail` | `create`、`read`、`gain`、`del` | `gm mail create 标题 内容` |
| `pub` | `multi_publish_test`、`game.StartGate`、`game.SyncPlayerInfo` | `gm pub multi_publish_test` |

说明：

- 物品 id 是 `Goods` 表的配置 id（如金币 `10100000`），不是物品实例 id
- 英雄相关命令的参数都是**英雄配置 id**（如 `2090101`），卡牌模型下英雄没有独立实例 id
- GM 命令要求账号 `privilege > 0`（`Account.Init` 默认给 3），权限不足会返回 `privilege not enough`
