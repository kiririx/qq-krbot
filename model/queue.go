package model

import (
	"container/list"
	"fmt"
)

type Queue struct {
	l []any
}

func f() {
	stack := list.New() // 新建栈
	stack.PushBack(1)   // 压栈
	stack.PushBack(2)
	stack.PushBack(3)
	fmt.Println(stack.Back().Value) // 输出栈顶元素
	stack.Remove(stack.Back())      // 弹栈

}
