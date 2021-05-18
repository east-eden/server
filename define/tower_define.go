package define

const (
	Tower_Type_Begin    int32 = iota
	Tower_Type_Melody   int32 = iota - 1 // 0 韵律之塔
	Tower_Type_Nature                    // 1 自然之塔
	Tower_Type_Civilize                  // 2 文明之塔
	Tower_Type_Destroy                   // 3 破灭之塔
	Tower_Type_General                   // 4 综合试炼
	Tower_Type_End
)

const TowerMaxFloor int32 = 100
