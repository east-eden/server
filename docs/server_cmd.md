# 服务器命

客户端按`F1`调出控制台后可输入命令

## 登录命令
```
login [user_id] [user_name]
```
>例：
login 1 dudu

- `user_id` 平台帐户id，测试时可随意输入一组数字
	
- `user_name` 平台帐户名，测试时可随意输入一组字符串

如果`user_id`不存在，则创建账号，并且创建角色，角色名暂时与平台账户名一致


## 断开连接
```
disconnect
```

## 玩家命令
```
player exp [num]
player level [num]
```
>例：
player exp 1000
player level 11

- `num` 经验或等级数值

## 物品命令
```
item add [type_id]
item del [item_id]
item use [item_id]
item query
```
>例：
item add 1
item del 296784428775045699
item use 296784428775045699
item query

- `type_id` 物品的type_id，在ItemConfig中定义
- `item_id` 物品的序列id，唯一，在物品生成后可以通过item query命令查询到`item_id`

## 英雄命令
```
hero add [type_id]
hero del [hero_id]
hero query
```
>例：
hero add 1
hero del 296784428775045699
hero query

- `type_id` 英雄的type_id，在HeroConfig中定义
- `hero_id` 英雄的序列id，唯一，在英雄生成后可以通过hero query命令查询到`hero_id`
