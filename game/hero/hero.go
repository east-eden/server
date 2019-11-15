package hero

type Hero interface {
	Init() error
	ID() int64
	Entry() *HeroEntry
}

var (
	DefaultHero defaultHero = newDefaultHero()
)

func NewHero() Hero {
	return newDefaultHero()
}
