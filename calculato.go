package main

import (
	"fmt"
	"math"
	"strconv"
)

var heapPriority map[rune]int = map[rune]int{'(': 0, '+': 1, '-': 1, '*': 2, '/': 2, '%': 2, '^': 3, '!': 4}
var priority map[rune]int = map[rune]int{')': 0, '+': 1, '-': 1, '*': 2, '/': 2, '%': 2, '^': 3, '!': 4, '(': 5}

type object interface{}

type stack struct {
	data []object
	top  int
}

func (s *stack) add(d object) {
	s.data = append(s.data, d)
	s.top++
}

func (s *stack) det() {
	if s.top >= 0 {
		s.data = s.data[:s.top]
		s.top--
	}
}

func (s *stack) read() object {
	if s.top >= 0 {
		return s.data[s.top]
	}
	return nil
}

func newStack() stack {
	return stack{top: -1}
}

// DocCalculato 计算器功能文档
var DocCalculato = &HelpDoc{
	Name:        "计算器",
	KeyWord:     []string{"计算器", "计算"},
	Example:     "计算器 17+4/2-8*(2^3*3+4-3!)%5",
	Description: "目前支持加减乘除开方求模阶乘"}

func calculato(msg []string, msgID int32, group, qq int64, try uint8) {

	if len(msg) == 0 { // 如果没有获取到参数
		try++         // 已尝试次数+1
		if try <= 3 { // 如果已尝试次数不超过3次
			if try == 1 {
				sendMsg(group, qq, "请输入算式\nq(≧▽≦q)")
			} else {
				sendMsg(group, qq, "没有算式怎么算,再发一次吧\n(。・ω・。)") // 发送提示消息
			}
			stagedSessionPool[msgID] = newStagedSession(group, qq, calculato, msg, try) // 添加新的会话到会话池
		} else {
			sendMsg(group, qq, "错误次数太多了哦,先看看使用说明吧\n(。・ω・。)") // 发送提示消息
		}
		return
	}

	expression := msg[0]
	operatioResult, err := calculation(expression)
	if err != nil {
		sendMsg(group, qq, fmt.Sprintln(err))
		return
	}
	sendMsg(group, qq, strconv.FormatFloat(operatioResult, 'G', 15, 64))
}

func operation(a, b float64, s rune) (float64, error) {
	switch s {
	case '+':
		return a + b, nil
	case '-':
		return a - b, nil
	case '*':
		return a * b, nil
	case '/':
		if b != 0 {
			return a / b, nil
		}
		return 0, fmt.Errorf("除以零错误")
	case '%':
		return math.Mod(a, b), nil
	case '^':
		return math.Pow(a, b), nil
	}
	return 0, nil
}

func factorial(a float64) (float64, error) {
	if a == 0 {
		return 1, nil
	}
	i, _ := math.Modf(a)
	if i == a {
		return math.Gamma(a + 1), nil
	}
	return math.Gamma(a), nil
}

func addNum(numS *[]rune, num *stack) error {
	if len(*numS) != 0 {
		i, err := strconv.ParseFloat(string(*numS), 10)
		if err != nil {
			return fmt.Errorf("\"%s\" 不是合法的数值", string(*numS))
		}
		num.add(i)
		*numS = []rune{}
		return nil
	}
	return fmt.Errorf("非法算式")
}

func calculation(input string) (output float64, err error) {
	var numS []rune
	var num stack = newStack()

	var signs stack = newStack()
	var flag bool = false
	num.add(0.0)
	for _, i := range input {
		switch i {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.': // 数字元素
			numS = append(numS, i)

		case '(': // 括号优先
			if len(numS) != 0 {
				return 0, fmt.Errorf("语法错误:意外的 %q", i)
			}
			signs.add(i)

		case '!': // 一元阶乘
			err = addNum(&numS, &num)
			if err != nil {
				return 0, err
			}
			n := num.read().(float64)
			num.det()
			o, err := factorial(n)
			if err != nil {
				return 0, err
			}
			num.add(o)
			flag = true

		case '-': // 判断负号
			if len(numS) == 0 && flag == false {
				numS = append(numS, i)
				break
			}
			fallthrough

		case '+', '*', '/', '^', '%': // 是二元运算符
			if flag {
				flag = false
			} else {
				err = addNum(&numS, &num)
				if err != nil {
					return 0, err
				}
			}

			for {
				ns := signs.read() // 栈顶运算符

				if ns != nil && heapPriority[ns.(rune)] >= priority[i] { // 栈不为空 且 栈顶运算符优先级 不低于 当前运算符

					a := num.read().(float64) // 取出 后运算值
					num.det()
					b := num.read().(float64) // 取出 前运算值
					num.det()
					signs.det()

					c, err := operation(b, a, ns.(rune)) // 运算
					if err != nil {
						return 0, err
					}

					num.add(c) // 储存结果

				} else {
					signs.add(i)
					break
				}
			}
		case ')':
			// fmt.Println(&numS, &num, "172 line")
			if flag {
				flag = false
			} else {
				err = addNum(&numS, &num)
				if err != nil {
					return 0, err
				}
			}
			ns := signs.read() // 栈顶运算符

			for ns != '(' {
				if ns == nil {
					return 0, fmt.Errorf("不匹配的')'")
				}
				a := num.read().(float64) // 取出 后运算值
				num.det()
				b := num.read().(float64) // 取出 前运算值
				num.det()
				signs.det()

				c, err := operation(b, a, ns.(rune)) // 运算
				if err != nil {
					return 0, err
				}

				num.add(c)        // 储存结果
				ns = signs.read() // 栈顶运算符
			}
			signs.det()
			flag = true

		default: // 其他字符
			return 0, fmt.Errorf("意外的 %q,\"%s\" 不是合法算式", i, input)
		}
	}

	addNum(&numS, &num)

	// fmt.Printf("运算数堆栈:%v  运算符堆栈:%q\n", num, signs)
	for { // 扫描完成,扫尾堆栈
		s := signs.read()         // 取出算子
		a := num.read().(float64) // 取出后项
		signs.det()
		num.det()
		if s == nil { // 算子为空,弹出栈顶值
			return a, nil
		} else if s == '(' {
			return 0, fmt.Errorf("不匹配的'('")
		}
		b := num.read().(float64) // 取出前项
		num.det()

		c, err := operation(b, a, s.(rune)) // 计算
		if err != nil {
			return 0, err
		}
		num.add(c) // 结果压入栈顶
	}
}
