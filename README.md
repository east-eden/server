[![Build Status](https://travis-ci.com/east-eden/server.svg?branch=master)](https://travis-ci.com/east-eden/server)
[![GoDoc](https://godoc.org/github.com/east-eden/server?status.svg)](https://godoc.org/github.com/east-eden/server)
[![Go Report Card](https://goreportcard.com/badge/github.com/east-eden/server)](https://goreportcard.com/report/github.com/east-eden/server)

# server
server is a game server with horizontally-scalable and high-available. It was powered by [go-micro](https://github.com/micro/go-micro) and running in docker container.

- [简体中文手册](docs/manual.md)
- [organize graph](docs/organize_graph.md)
- [tcp protocol](docs/tcp_protocol.md)

## Requirement
- **MongoDB**

## Getting Started
- **Download** - git clone this repo and cd in its root path

- **Start MongoDB** - running in `docker-compose`:
```
docker-compose run --service-ports -d mongo
```

- **Start Gate** - cd to `apps/gate`, run following command:
```
go run main.go
```

- **Start Game** - open another terminal session, cd to `apps/game`, run following command:
```
go run main.go
```

- **Start Combat** - open another terminal session, cd to `apps/combat`, run following command:
```
go run main.go
```

- **Start Client** - open another terminal session, cd to `apps/client`, run following command:
```
go run main.go
```
now you can communicate with server using (up down left right enter):
![text mud](https://raw.githubusercontent.com/east-eden/server/master/docs/text_mud.jpg)

## Using store to save object in cache and database
- first add a new store info
```golang
func init() {
    // add store info
    store.GetStore().AddStoreInfo(define.StoreType_Account, "account", "_id")
    store.GetStore().AddStoreInfo(define.StoreType_Player, "player", "_id")
    store.GetStore().AddStoreInfo(define.StoreType_Item, "item", "_id")
    store.GetStore().AddStoreInfo(define.StoreType_Hero, "hero", "_id")
    store.GetStore().AddStoreInfo(define.StoreType_Token, "token", "_id")
}

```

- load single object example

```golang
func (m *TokenManager) LoadAll() {
	err := store.GetStore().FindOne(context.Background(), define.StoreType_Token, m.owner.GetID(), m)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}
}
```

- load array example

```golang
func (m *HeroManager) LoadAll() error {
	res, err := store.GetStore().FindAll(context.Background(), define.StoreType_Hero, "owner_id", m.owner.ID)
	if errors.Is(err, store.ErrNoResult) {
		return nil
	}

	if !utils.ErrCheck(err, "FindAll failed when HeroManager.LoadAll", m.owner.ID) {
		return err
	}

	for _, v := range res {
		vv := v.([]byte)
		h := hero.NewHero()
		err := json.Unmarshal(vv, h)
		if !utils.ErrCheck(err, "Unmarshal failed when HeroManager.LoadAll") {
			continue
		}

		if err := m.initLoadedHero(h); err != nil {
			return fmt.Errorf("HeroManager LoadAll: %w", err)
		}
	}

	return nil
}
```

- save example

```golang
func (m *HeroManager) AddHeroByTypeId(typeId int32) *hero.Hero {
	heroEntry, ok := auto.GetHeroEntry(typeId)
	if !ok {
		log.Warn().Int32("type_id", typeId).Msg("GetHeroEntry failed")
		return nil
	}

	// 重复获得卡牌，转换为对应碎片
	_, ok = m.heroTypeSet[typeId]
	if ok {
		m.owner.FragmentManager().Inc(typeId, heroEntry.FragmentTransform)
		return nil
	}

	h := m.createEntryHero(heroEntry)
	if h == nil {
		log.Warn().Int32("type_id", typeId).Msg("createEntryHero failed")
		return nil
	}

	err := store.GetStore().UpdateOne(context.Background(), define.StoreType_Hero, h.Id, h)
	if !utils.ErrCheck(err, "UpdateOne failed when AddHeroByTypeID", typeId, m.owner.ID) {
		m.delHero(h)
		return nil
	}

	m.SendHeroUpdate(h)

	// prometheus ops
	prom.OpsCreateHeroCounter.Inc()

	return h
}
```

- save several fields example

```golang
func makeTokenKey(tp int32) string {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)

	_, _ = b.WriteString("tokens.")
	_, _ = b.WriteString(cast.ToString(tp))

	return b.String()
}

func (m *TokenManager) save(tp int32) error {
	fields := map[string]interface{}{
		makeTokenKey(tp): m.Tokens[tp],
	}
	return store.GetStore().UpdateFields(context.Background(), define.StoreType_Token, m.owner.ID, fields)
}
```



## License
server is MIT licensed.

