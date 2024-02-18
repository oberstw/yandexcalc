package orch

import (
	"net/http"
	"fmt"
	"orch/handlers"
    "time"
    "context"
    "strconv"
)

type Stack []string

type Agentreq struct {
	A     float64 `json:"a"`
	B     float64 `json:"b"`
	Sign    string  `json:"sign"`
	Timeout int     `json:"timeout"`
}

type Agentout struct {
	Result float64 `json:"result"`
	Err string  `json:"err"`
}

func (st *Stack) IsEmpty() bool {
    return len(*st) == 0
}

func (st *Stack) Push(str string) {
    *st = append(*st, str) 
}

func (st *Stack) Pop() bool {
    var x
    x, *st = *st[len(*st)-1], *st[:len(*st)-1]
    return x
}

func (st *Stack) Top() string {
    if st.IsEmpty() {
        return ""
    } else {
        index := len(*st) - 1   
        element := (*st)[index] 
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

func InfixToPostfix(infix string) string {
    var sta Stack
    var postfix []string
    for _, char := range infix {
        opchar := string(char)
        if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
            postfix = append(postfix, opchar)
        } else if char == '(' {
            sta.Push(opchar)
        } else if char == ')' {
            for sta.Top() != "(" {
                postfix = append(postfix, sta.Top())
                sta.Pop()
            }
            sta.Pop()
        } else {
            for !sta.IsEmpty() && prec(opchar) <= prec(sta.Top()) {
                postfix = append(postfix, sta.Top())
                sta.Pop()
            }
            sta.Push(opchar)
        }
    }
    for !sta.IsEmpty() {
        postfix = append(postfix, sta.Top())
        sta.Pop()
    }
    return postfix
}

func sanitize(line string) string{
	return strings.ReplaceAll(line, " ", "")
}

func Calculate(arr []string) float64, error {
    var stack Stack
	for _, i := range arr {
        if num, err := strconv.ParseFloat(i, 64); err == nil {
			stack.Push(num)
		} else {
			if stack.Len() < 2 {
				return 0, fmt.Errorf("invalid postfix expression")
			}
			b := stack.Pop()
			a := stack.Pop()

			topush, err := AgentCalc(a, b, i)
			if err != nil{
				return 0, err
			}
			stack.Push(topush)
		}
	}
}

func AgentCalc(a, b float64, sign string) (float64, error) {
    url := "http://127.0.0.1:8080/exp"
	client := &http.Client{}
	s := Agentreq{a, b, sign, 1000}
    data, _ := json.Marshal(s)
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1000) * time.Millisecond)
	defer cancel()

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, bytes.NewBuffer(jsdata))
    if err != nil{
		return 0, fmt.Errorf("bad agent request")
	}
    resp, err := client.Do(req)
    if err != nil {
        if err == context.DeadlineExceeded{
			return 0, fmt.Errorf("timeout exceeded")
		} else {
			return 0, fmt.Errorf("request failed")
		}
    }
    defer resp.Body.Close()
    var data Agentout
    decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if resp.StatusCode != http.StatusOK{
		return 0, fmt.Errorf(data.Err)
	}
    return data.Result, nil
}