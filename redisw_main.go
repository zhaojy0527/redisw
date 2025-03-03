package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2"
)

type RedisServer struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

var (
	configFile = flag.String("config", getDefaultConfigPath(), "path to config file")
)

func getDefaultConfigPath() string {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "./redisw_config.yml"
    }

    // 首先检查家目录
    homeYml := filepath.Join(homeDir, "redisw_config.yml")
    homeYaml := filepath.Join(homeDir, "redisw_config.yaml")

    if _, err := os.Stat(homeYml); err == nil {
        return homeYml
    }
    if _, err := os.Stat(homeYaml); err == nil {
        return homeYaml
    }

    // 其次检查 .config 目录
    configDir := filepath.Join(homeDir, ".config", "redisw")
    if err := os.MkdirAll(configDir, 0755); err != nil {
        return "./redisw_config.yml"
    }

    // 检查 ~/.config/redisw 目录下的配置文件
    ymlPath := filepath.Join(configDir, "redisw_config.yml")
    yamlPath := filepath.Join(configDir, "redisw_config.yaml")

    // 如果配置文件已存在，直接返回
    if _, err := os.Stat(ymlPath); err == nil {
        return ymlPath
    }
    if _, err := os.Stat(yamlPath); err == nil {
        return yamlPath
    }

    // 尝试从项目目录复制配置文件
    projectYml := "./redisw_config.yml"
    projectYaml := "./redisw_config.yaml"

    // 检查并复制项目目录的配置文件
    if _, err := os.Stat(projectYml); err == nil {
        copyFile(projectYml, ymlPath)
        return ymlPath
    }
    if _, err := os.Stat(projectYaml); err == nil {
        copyFile(projectYaml, ymlPath)
        return ymlPath
    }

    // 如果都不存在，创建默认配置
    defaultConfig := []RedisServer{
        {
            Name:     "localhost",
            Host:     "127.0.0.1",
            Port:     6379,
            Password: "",
        },
    }

    file, err := os.Create(ymlPath)
    if err == nil {
        encoder := yaml.NewEncoder(file)
        encoder.Encode(defaultConfig)
        file.Close()
    }

    return ymlPath
}

// 添加复制文件的辅助函数
func copyFile(src, dst string) error {
    input, err := os.ReadFile(src)
    if err != nil {
        return err
    }

    err = os.WriteFile(dst, input, 0644)
    if err != nil {
        return err
    }

    return nil
}

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
		Size:  len(servers),            // 设置显示的项目数量为所有服务器
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
	cmd := exec.Command("redis-cli", "-h", server.Host, "-p", fmt.Sprintf("%d", server.Port), "-c")
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
