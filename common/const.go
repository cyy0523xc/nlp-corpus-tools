package common

var (
	// 全局测试状态
	Debug bool
	Skip  int // 跳过前面的若干行记录
	Limit int // 只读取若干条有效的记录
)
