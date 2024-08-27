package capture

import "fmt"

// 定义一个自定义错误类型
type E2EmptyConfigError struct {
}

// 实现Error() string方法
func (ce E2EmptyConfigError) Error() string {
	return fmt.Sprintf("not init capture config")
}
