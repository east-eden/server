package hero

type defaultHero struct {
	id    int64
	entry *HeroEntry
}

func newDefaultHero() Hero {
	return &defaultHero{}
}

func (h *defaultHero) Init() error {
	h.id = 1
	return nil
}

func (h *defaultHero) ID() int64 {
	return h.id
}

func (h *defaultHero) Entry() *HeroEntry {
	return h.entry
}
