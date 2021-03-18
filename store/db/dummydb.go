package db

type DummyDB struct {
}

func NewDummyDB() DB {
	m := &DummyDB{}

	return m
}

// migrate collection
func (m *DummyDB) MigrateTable(name string, indexNames ...string) error {
	return nil
}

func (m *DummyDB) LoadObject(tblName, key string, value interface{}, x interface{}) error {
	return nil
}

func (m *DummyDB) LoadArray(tblName, key string, value interface{}, x interface{}) error {
	return nil
}

func (m *DummyDB) SaveObject(tblName string, k interface{}, x interface{}) error {
	return nil
}

func (m *DummyDB) SaveFields(tblName string, k interface{}, fields map[string]interface{}) error {
	return nil
}

func (m *DummyDB) DeleteObject(tblName string, k interface{}) error {
	return nil
}

func (m *DummyDB) DeleteFields(tblName string, k interface{}, fieldsName []string) error {
	return nil
}

func (m *DummyDB) Exit() {
}
