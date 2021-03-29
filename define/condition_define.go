package define

const (
	Condition_Type_Begin int32 = iota
	Condition_Type_And   int32 = iota - 1 // 需满足所有子条件
	Condition_Type_Or                     // 满足其中一个子条件

	Condition_Type_End
)

const (
	Condition_SubType_Begin             int32 = iota
	Condition_SubType_TeamLevel_Achieve int32 = iota - 1 // 0 队伍等级达到**级

	Condition_SubType_End
)
