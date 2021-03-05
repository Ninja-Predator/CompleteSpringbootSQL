package main

import (
	"bufio"
	"fmt"
	"github.com/atotto/clipboard"
	"os"
	"strconv"
	"strings"
)

func main() {
	for true {
		formatSql()
		fmt.Print("回车再来一次？")
		bufio.NewReader(os.Stdin).ReadLine()
	}
}

func formatSql() {
	p := "Preparing: "
	ps := "Parameters: "
	var PSQL []string
	var Parameters []string
	cBoard, _ := clipboard.ReadAll()
	lines := strings.Split(cBoard, "\n")
	for _, line := range lines {
		if strings.Contains(line, p) {
			PSQLs := strings.Split(line, p)
			PSQL = strings.Split(PSQLs[len(PSQLs)-1], "?")
		} else if strings.Contains(line, ps) {
			Parameters1 := strings.Split(line, ps)
			Parameters = strings.Split(Parameters1[len(Parameters1)-1], "), ")
		}
	}
	completeSql(PSQL, Parameters)
}

func completeSql(PSQL []string, Parameters []string) {
	if len(PSQL) != len(Parameters)+1 {
		if len(PSQL) == 0 || len(Parameters) == 0 {
			fmt.Println("你没复制东西吧")
			return
		}
		fmt.Println("拆分长度出错")
		fmt.Println("SQL长度: ", len(PSQL))
		fmt.Println("参数长度: ", len(Parameters))
		for i, Para := range Parameters {
			if strings.Contains(Para, "null, ") {
				fmt.Println("我大意了啊，怎么在第" + strconv.Itoa(i+1) + "位的参数有个null，塞回去我再看看啊")
				temp := strings.Split(Para, "null, ")[1]
				var ss []string
				if i != 0 {
					ss = append(ss, Parameters[:i-1]...)
				}
				ss = append(ss, "null(null)", temp)
				ss = append(ss, Parameters[i+1:]...)
				completeSql(PSQL, ss)
				break
			}
		}
	} else {
		var finalSql strings.Builder
		finalSql.WriteString(PSQL[0])
		for index, ParameterW := range Parameters {
			Parameter := ParameterW[:strings.LastIndexAny(ParameterW, "(")]
			St := strings.Split(ParameterW, "(")[len(strings.Split(ParameterW, "("))-1]
			if strings.Contains(St, "String") {
				Parameter = "'" + Parameter + "'"
			}
			finalSql.WriteString(Parameter)
			finalSql.WriteString(PSQL[index+1])
		}
		str := strip(finalSql.String(), "\n\r\t ")
		str += ";"
		clipboard.WriteAll(str)
		fmt.Println(str)
		fmt.Println("拼好的sql放到剪贴板了，别手动复制了")
	}
}

func strip(s_ string, chars_ string) string {
	s, chars := []rune(s_), []rune(chars_)
	length := len(s)
	max := len(s) - 1
	l, r := true, true //标记当左端或者右端找到正常字符后就停止继续寻找
	start, end := 0, max
	tmpEnd := 0
	charset := make(map[rune]bool) //创建字符集，也就是唯一的字符，方便后面判断是否存在
	for i := 0; i < len(chars); i++ {
		charset[chars[i]] = true
	}
	for i := 0; i < length; i++ {
		if _, exist := charset[s[i]]; l && !exist {
			start = i
			l = false
		}
		tmpEnd = max - i
		if _, exist := charset[s[tmpEnd]]; r && !exist {
			end = tmpEnd
			r = false
		}
		if !l && !r {
			break
		}
	}
	if l && r { // 如果左端和右端都没找到正常字符，那么表示该字符串没有正常字符
		return ""
	}
	return string(s[start : end+1])
}
