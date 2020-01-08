package define

// Global mysql table global
type TableGlobal struct {
	ID        int `gorm:"type:int(10);primary_key;column:id;default:0;not null;unsigned" bson:"_id"`
	TimeStamp int `gorm:"type:int(10);column:time_stamp;default:0;not null" bson:"timestamp"`
}

// TableName set global table name to be `global`
func (TableGlobal) TableName() string {
	return "global"
}

// gate mysql table
type TableGate struct {
	ID        int `gorm:"type:int(10);primary_key;column:id;default:0;not null;unsigned" bson:"_id"`
	TimeStamp int `gorm:"type:int(10);column:time_stamp;default:0;not null" bson:"timestamp"`
}

// TableName set gate table name to be `gate`
func (TableGate) TableName() string {
	return "gate"
}
