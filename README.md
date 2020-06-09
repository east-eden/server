[![Build Status](https://travis-ci.com/yokaiio/yokai_server.svg?branch=master)](https://travis-ci.com/yokaiio/yokai_server)
[![GoDoc](https://godoc.org/github.com/yokaiio/yokai_server?status.svg)](https://godoc.org/github.com/yokaiio/yokai_server)
[![Go Report Card](https://goreportcard.com/badge/github.com/yokaiio/yokai_server)](https://goreportcard.com/report/github.com/yokaiio/yokai_server)

# yokai_server
yokai_server is a game server with horizontally-scalable and high-available. It was powered by [go-micro](https://github.com/micro/go-micro) and running in docker container.

- [organize graph](docs/organize_graph.md)
- [tcp protocol](docs/tcp_protocol.md)

## Requirement
- **MongoDB**
- **Redis**

## Getting Started
- **Download** - git clone this repo and cd in its root path

- **Start MongoDB** - running in `docker-compose`:
```
docker-compose run --service-ports -d mongo
```

- **Start Redis** - running in `docker-compose`:
```
docker-compose run --service-ports -d redis
```

- **Start Gate** - cd to `apps/gate`, run following command:
```
go run main.go plugin.go
```

- **Start Game** - open another terminal session, cd to `apps/game`, run following command:
```
go run main.go plugin.go
```

- **Start Combat** - open another terminal session, cd to `apps/combat`, run following command:
```
go run main.go plugin.go
```

- **Start Client** - open another terminal session, cd to `apps/client`, run following command:
```
go run main.go plugin.go
```
now you can communicate with server using (up down left right enter):
![text mud](https://raw.githubusercontent.com/yokaiio/yokai_server/master/docs/text_mud.jpg)

## Using store to save object in cache and database

- load example

```golang
heroList, err := store.GetStore().LoadArrayFromCacheAndDB(store.StoreType_Hero, "owner_id", m.owner.GetID(), hero.GetHeroPool())
	if err != nil {
		logger.Error("load hero manager failed:", err)
	}
```

- save full object example

```golang
h := hero.NewHero(
    hero.Id(id),
    hero.OwnerId(m.owner.GetID()),
    hero.OwnerType(m.owner.GetType()),
    hero.Entry(entry),
    hero.TypeId(entry.ID),
)

//h.GetAttManager().SetBaseAttId(entry.AttID)
//m.mapHero[h.GetOptions().Id] = h
store.GetStore().SaveObjectToCacheAndDB(store.StoreType_Hero, h)
```

-- save several fields example

```golang
func (m *TokenManager) save() error {
	fields := map[string]interface{}{
		"tokens":           m.Tokens,
		"serialize_tokens": m.SerializeTokens,
	}
	store.GetStore().SaveFieldsToCacheAndDB(store.StoreType_Token, m, fields)

	return nil
}
```

## License
yokai_server is MIT licensed.

