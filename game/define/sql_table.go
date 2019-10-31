package define

// Global mysql table global
type TableGlobal struct {
	ID        uint `gorm:"type:int(10);primary_key;column:id;default:0;not null;unsigned"`
	TimeStamp int  `gorm:"type:int(10);column:time_stamp;default:0;not null"`
}

// TableName set global table name to be `global`
func (TableGlobal) TableName() string {
	return "global"
}
