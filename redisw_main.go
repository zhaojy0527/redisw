package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2" // 用于解析 YAML 文件
)

type RedisServer struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

var (
	configFile = flag.String("config", "./redisw_config.yml", "path to config file")
)

func main() {
	flag.Parse()
	// 从 YAML 文件加载配置
	servers := loadConfig(*configFile)

	for {
		// 提示用户选择 Redis 服务器
		selectedServer := chooseServer(servers)
		if selectedServer == nil {
			return
		}

		// 使用 redis-cli 连接到选定的 Redis 服务器
		connectToRedis(selectedServer)
	}
}

func loadConfig(filePath string) []RedisServer {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return nil
	}
	defer file.Close()

	var servers []RedisServer
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&servers)
	if err != nil {
		fmt.Println("Error decoding config file:", err)
		return nil
	}

	return servers
}

func chooseServer(servers []RedisServer) *RedisServer {
	// 创建一个 promptui.Select 对象
	prompt := promptui.Select{
		Label: "Select Redis Server",
		Items: getServerNames(servers), // 使用辅助函数获取服务器名称
		Size:  len(servers), // 设置显示的项目数量为所有服务器
		Searcher: func(input string, index int) bool {
			server := servers[index]
			name := strings.ReplaceAll(strings.ToLower(server.Name), " ", "")
			return strings.Contains(name, strings.ToLower(input))
		},
		Templates: &promptui.SelectTemplates{
			Label:    "✨ {{ . | green}}",
			Selected: `{{ "➤ " | green }}{{ . | faint }}`, // 显示选中的服务
			Active:   `{{ "➤ " | green }}{{ . | green }}`, // 高亮显示选中的服务
		},
	}

	// 运行选择提示
	index, _, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return nil
	}

	// 返回选定的服务器
	return &servers[index]
}

func connectToRedis(server *RedisServer) {
	// 构建 redis-cli 命令
	cmd := exec.Command("redis-cli", "-h", server.Host, "-p", fmt.Sprintf("%d", server.Port),"-c")
	if server.Password != "" {
		cmd.Args = append(cmd.Args, "-a", server.Password)
	}

	// 设置命令的标准输入输出
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 启动 redis-cli
	err := cmd.Start()
	if err != nil {
		fmt.Println("Failed to start redis-cli:", err)
		return // 这里不退出程序，允许返回选择服务器
	}

	// 等待 redis-cli 完成
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Failed to connect to Redis using redis-cli:", err)
		return // 这里不退出程序，允许返回选择服务器
	}
}

// 辅助函数：获取服务器名称
func getServerNames(servers []RedisServer) []string {
	names := make([]string, len(servers))
	for i, server := range servers {
		names[i] = server.Name
	}
	return names
}