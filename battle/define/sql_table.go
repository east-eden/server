package define

// battle mysql table
type TableBattle struct {
	ID        int `gorm:"type:int(10);primary_key;column:id;default:0;not null;unsigned"`
	TimeStamp int `gorm:"type:int(10);column:time_stamp;default:0;not null"`
}

// TableName set battle table name to be `battle`
func (TableBattle) TableName() string {
	return "battle"
}
