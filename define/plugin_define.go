package define

const (
	Plugin_Player = iota
	Plugin_Hero
	Plugin_Item

	Plugin_End
)

type PluginObj interface {
	GetType() int32
	GetID() int64
	GetLevel() int32
}
