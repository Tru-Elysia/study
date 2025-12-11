package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// 定义用户消息结构体
type userMessages struct {
	username string
	messages string
}

// 定义用户消息记录列表
var mesLog []userMessages

// 定义在线用户列表
var onlineUsers []string

var mu sync.RWMutex

// 处理用户登录
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// 仅处理POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}
	// 解析表单数据
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "表单解析失败："+err.Error(), http.StatusBadRequest)
		return
	}
	// 处理用户登录
	mu.Lock()
	defer mu.Unlock()
	username := strings.TrimSpace(r.Form.Get("username"))
	if username == "" {
		http.Error(w, "用户名不能为空", http.StatusBadRequest)
		return
	}

	// 添加用户到在线用户列表                                                    

	// 回复登录成功
	fmt.Fprintf(w, "用户 %s 登录成功！当前在线：%d 人", username, len(onlineUsers))

	// 广播系统通知
	joinMsg := userMessages{
		username: "【系统通知】",
		messages: fmt.Sprintf("%s 加入聊天室，当前在线 %d 人", username, len(onlineUsers)),
	}
	mesLog = append(mesLog, joinMsg)
}

// 处理用户登出
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// 仅处理POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}
	// 解析表单数据
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "表单解析失败："+err.Error(), http.StatusBadRequest)
		return
	}

	// 处理用户登出
	mu.Lock()
	defer mu.Unlock()
	username := strings.TrimSpace(r.Form.Get("username"))
	if username == "" {
		http.Error(w, "用户名不能为空", http.StatusBadRequest)
		return
	}

	// 检查用户是否存在
	index := -1
	for i, user := range onlineUsers {
		if user == username {
			index = i
			break
		}
	}
	if index == -1 {
		http.Error(w, "用户未在线", http.StatusBadRequest)
		return
	}
	// 从在线用户列表移除用户
	onlineUsers = append(onlineUsers[:index], onlineUsers[index+1:]...)
	// 回复登出成功
	fmt.Fprintf(w, "已登出！当前在线：%d 人", len(onlineUsers))
	// 广播系统通知
	logoutMsg := userMessages{
		username: "【系统通知】",
		messages: fmt.Sprintf("%s 离开聊天室，当前在线 %d 人", username, len(onlineUsers)),
	}
	mesLog = append(mesLog, logoutMsg)
}

// 处理用户发送消息
func sendHandler(w http.ResponseWriter, r *http.Request) {
	// 仅处理POST请求
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}
	// 解析表单数据
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "表单解析失败："+err.Error(), http.StatusBadRequest)
		return
	}
	// 处理用户发送消息
	mu.RLock()
	username := strings.TrimSpace(r.Form.Get("username"))
	message := strings.TrimSpace(r.Form.Get("message"))

	// 检查用户名和消息内容是否为空
	if username == "" || message == "" {
		http.Error(w, "用户名或消息不能为空", http.StatusBadRequest)
		return
	}

	// 检查用户是否存在
	ok := false
	for _, user := range onlineUsers {
		if user == username {
			ok = true
			break
		}
	}
	if !ok {
		http.Error(w, "用户未在线", http.StatusBadRequest)
		return
	}
	mu.RUnlock()

	// 记录用户消息
	mu.Lock()
	defer mu.Unlock()
	userMsg := userMessages{
		username: username,
		messages: message,
	}
	mesLog = append(mesLog, userMsg)
	// 回复发送成功
	fmt.Fprintf(w, "消息发送成功！")
}

// 获取消息记录
func messageHandler(w http.ResponseWriter, r *http.Request) {
	// 仅处理GET请求
	if r.Method != http.MethodGet {
		http.Error(w, "仅支持GET请求", http.StatusMethodNotAllowed)
		return
	}
	mu.RLock()
	defer mu.RUnlock()

	// 输出历史消息
	for _, msg := range mesLog {
		fmt.Fprintf(w, "%s: %s\n", msg.username, msg.messages)
	}
}

// 处理在线用户列表请求
func onlineuserHandler(w http.ResponseWriter, r *http.Request) {
	// 仅处理GET请求
	if r.Method != http.MethodGet {
		http.Error(w, "仅支持GET请求", http.StatusMethodNotAllowed)
		return
	}
	mu.RLock()
	defer mu.RUnlock()
	// 输出在线用户列表
	fmt.Fprintf(w, "当前共 %d 位用户在线：\n", len(onlineUsers))
	for _, user := range onlineUsers {
		fmt.Fprintf(w, "%s\n", user)
	}
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/message", messageHandler)
	http.HandleFunc("/onlineuser", onlineuserHandler)

	// 启动服务器
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("服务器启动失败，错误：", err)
		return
	}
	
}
