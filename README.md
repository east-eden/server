[![Build Status](https://travis-ci.com/yokaiio/yokai_server.svg?branch=master)](https://travis-ci.com/yokaiio/yokai_server)
[![GoDoc](https://godoc.org/github.com/yokaiio/yokai_server?status.svg)](https://godoc.org/github.com/yokaiio/yokai_server)
[![Go Report Card](https://goreportcard.com/badge/github.com/yokaiio/yokai_server)](https://goreportcard.com/report/github.com/yokaiio/yokai_server)

# yokai_server
yokai_server is a game server with horizontally-scalable and high-available. It was powered by [go-micro](https://github.com/micro/go-micro) and running in docker container.

- [organize graph](docs/organize_graph.md)
- [tcp protocol](docs/tcp_protocol.md)

## Requirement
- **MongoDB**
- **Redis-json module**

## Getting Started
- **Download** - git clone this repo and cd in its root path

- **Start MongoDB** - running in `docker-compose`:
```
docker-compose run --service-ports -d mongo
```

- **Start Redis with json module** - running in `docker-compose`:
```
docker-compose run --service-ports -d rejson
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
![text mud](https://raw.githubusercontent.com/yokaiio/yokai_server/master/docs/text_mud.jpg)

## Using store to save object in cache and database
- first add a new store info
```golang
func init() {
    // add store info
    store.GetStore().AddStoreInfo(define.StoreType_Account, "account", "_id", "")
    store.GetStore().AddStoreInfo(define.StoreType_Player, "player", "_id", "")
    store.GetStore().AddStoreInfo(define.StoreType_Item, "item", "_id", "owner_id")
    store.GetStore().AddStoreInfo(define.StoreType_Hero, "hero", "_id", "owner_id")
    store.GetStore().AddStoreInfo(define.StoreType_Rune, "rune", "_id", "owner_id")
    store.GetStore().AddStoreInfo(define.StoreType_Token, "token", "_id", "owner_id")
}

```

- load single object example

```golang
func (m *TokenManager) LoadAll() {
	err := store.GetStore().LoadObject(define.StoreType_Token, m.owner.GetID(), m)
	if err != nil {
		store.GetStore().SaveObject(define.StoreType_Token, m)
		return
	}

}
```

- load array example

```golang
func (m *HeroManager) LoadAll() {
	heroList, err := store.GetStore().LoadArray(define.StoreType_Hero, m.owner.GetID(), hero.GetHeroPool())
	if err != nil {
		logger.Error("load hero manager failed:", err)
	}

	for _, i := range heroList {
		err := m.initLoadedHero(i.(hero.Hero))
		if err != nil {
			logger.Error("load hero failed:", err)
		}
	}
}
```

- save example

```golang
func (m *HeroManager) createEntryHero(entry *define.HeroEntry) hero.Hero {
    /*
	if entry == nil {
		logger.Error("newEntryHero with nil HeroEntry")
		return nil
	}

	id, err := utils.NextID(define.SnowFlake_Hero)
	if err != nil {
		logger.Error(err)
		return nil
	}
    */

	h := hero.NewHero(
		hero.Id(id),
		hero.OwnerId(m.owner.GetID()),
		hero.OwnerType(m.owner.GetType()),
		hero.Entry(entry),
		hero.TypeId(entry.ID),
	)

	//h.GetAttManager().SetBaseAttId(entry.AttID)
	//m.mapHero[h.GetOptions().Id] = h
	store.GetStore().SaveObject(define.StoreType_Hero, h)

	//h.GetAttManager().CalcAtt()

	return h
}
```

- save several fields example

```golang
func (m *TokenManager) save() error {
	fields := map[string]interface{}{
		"tokens":           m.Tokens,
		"serialize_tokens": m.SerializeTokens,
	}
	store.GetStore().SaveFields(define.StoreType_Token, m, fields)

	return nil
}
```

- full save and load example useing `lru cache` and `sync.Pool`

```golang
func (am *AccountManager) GetLitePlayer(playerId int64) (player.LitePlayer, error) {
	am.RLock()
	defer am.RUnlock()

	if lp, ok := am.litePlayerCache.Get(playerId); ok {
		return *(lp.(*player.LitePlayer)), nil
	}

	lp := am.litePlayerPool.Get().(*player.LitePlayer)
	err := store.GetStore().LoadObject(define.StoreType_LitePlayer, playerId, lp)
	if err == nil {
		am.litePlayerCache.Add(lp.ID, lp)
		return *lp, nil
	}

	am.litePlayerPool.Put(lp)
	return *(player.NewLitePlayer().(*player.LitePlayer)), err
}
```

## License
yokai_server is MIT licensed.

