package math

import (
	"net/http"
	"fmt"
    "time"
    "encoding/json"
    "context"
    "bytes"
    "strconv"
)

type Stack struct {
    Sli []string
}

type Stackint struct {
    Sl []float64
} 

var Time = make(map[string]int)

type Agentreq struct {
	A     float64 `json:"a"`
	B     float64 `json:"b"`
	Sign    string  `json:"sign"`
    Timeout int `json:"timeout"`
}

type Agentout struct {
	Result float64 `json:"result"`
	Err string  `json:"err"`
}

func (st *Stack) IsEmpty() bool {
    return len(st.Sli) == 0
}

func (st *Stack) Push(str string) {
    st.Sli = append(st.Sli, str) 
}

func (st *Stack) PopStack() string {
    var x string
    x, st.Sli = st.Sli[len(st.Sli)-1], st.Sli[:len(st.Sli)-1]
    return x
}

func (st *Stack) Top() string {
    if st.IsEmpty() {
        return ""
    } else {
        index := len(st.Sli) - 1   
        element := (st.Sli)[index] 
        return element
    }
}

func (st *Stackint) Push(i float64) {
    st.Sl = append(st.Sl, i) 
}

func (st *Stackint) IsEmpty() bool {
    return len(st.Sl) == 0
}

func (st *Stackint) Pop() float64 {
    var x float64
    x, st.Sl = st.Sl[len(st.Sl)-1], st.Sl[:len(st.Sl)-1]
    return x
}

func (st *Stackint) Len() float64 {
    return float64(len(st.Sl))
}

func (st *Stackint) Top() float64 {
    if st.IsEmpty() {
        return 0
    } else {
        index := len(st.Sl) - 1   
        element := (st.Sl)[index] 
        return element
    }
}

func prec(s string) int {
    if s == "^" {
        return 3
    } else if (s == "/") || (s == "*") {
        return 2
    } else if (s == "+") || (s == "-") {
        return 1
    } else {
        return -1
    }
}

func Oper(r rune) bool {
	return (r == 42 || r == 43 || r == 45 || r == 47)
}

func InfixToPostfix(infix string) []string {
    var stack Stack
    var postfix []string
    for _, char := range infix {
        opchar := string(char)
        if char >= '0' && char <= '9' {
            postfix = append(postfix, opchar)
            fmt.Println("Appended to postfix one op:" , postfix, opchar)

        } else if char == '(' {
            stack.Push(opchar)

        } else if char == ')' {
            for stack.Top() != "(" {
                postfix = append(postfix, stack.Top())
                stack.PopStack()
            }
            stack.PopStack()

        } else if Oper(char) {
            fmt.Println("Oper found: ", char)
            for !stack.IsEmpty() && prec(opchar) <= prec(stack.Top()) {
                fmt.Println("opchar", opchar)
                postfix = append(postfix, stack.Top())
                stack.PopStack()
            }
            fmt.Println("opchar - not in loop", opchar)
            stack.Push(opchar)
        } else {
            return postfix 
        }
        fmt.Println("loop stack:", stack)
    }
    fmt.Println(stack)
    for !stack.IsEmpty() {
        g := stack.Top()
        stack.PopStack()
        postfix = append(postfix, g)
    }
    fmt.Println("ItoP done", postfix)
    return postfix
}

func Calculate(arr []string) (float64, error) {
    var stack Stackint
    fmt.Println(arr)
	for _, i := range arr {
        if num, err := strconv.ParseFloat(i, 64); err == nil {
            stack.Push(num)
            fmt.Println(stack)
            fmt.Println(num)
		} else {
			if stack.Len() < 2 {
				return 0, fmt.Errorf("invalid postfix expression")
			}
			b := stack.Pop()
			a := stack.Pop()
            fmt.Println(stack)
            fmt.Println("Calculating: ", a, i, b)
			topush, err := AgentCalc(a, b, i)
            fmt.Println(err)
			if err != nil{
                fmt.Println("errored")
				return 0, err
			}
			stack.Push(topush)
		}
	}
    fmt.Println(stack.Sl[0])
    return stack.Sl[0], nil
}

func AgentCalc(a, b float64, sign string) (float64, error) {
    url := "http://127.0.0.1:8080/exp"
	client := &http.Client{}
	s := Agentreq{a, b, sign, Time[sign]}
    data, _ := json.Marshal(s)
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.Timeout + 100) * time.Millisecond)
	defer cancel()

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, bytes.NewBuffer(data))
    if err != nil{
		return 0, fmt.Errorf("bad agent request")
	}
    resp, err := client.Do(req)
    if err != nil {
        if err == context.DeadlineExceeded{
			return 0, fmt.Errorf("timeout exceeded")
		} else {
			return 0, err
		}
    }
    defer resp.Body.Close()
    var dat Agentout
    decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&dat)
    fmt.Println(dat.Result)
	if resp.StatusCode != http.StatusOK{
		return 0, fmt.Errorf(dat.Err)
	}
    return dat.Result, nil
}