package config

// 模式常量
const (
	Development ModeType = "development"
	Test        ModeType = "test"
	Production  ModeType = "production"
)

// ModeType 模式类型
type ModeType string

// IsValid 判定模式是否可用
func (m ModeType) IsValid() bool {
	switch m {
	case Development, Test, Production:
		return true

	}

	return false
}
