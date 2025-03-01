package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDefaultConfigPath(t *testing.T) {
	// 获取测试用的临时主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	// 测试配置目录
	configDir := filepath.Join(homeDir, ".config", "redisw")
	ymlPath := filepath.Join(configDir, "redisw_config.yml")
	yamlPath := filepath.Join(configDir, "redisw_config.yaml")

	// 清理测试环境
	defer func() {
		os.Remove(ymlPath)
		os.Remove(yamlPath)
		os.Remove(configDir)
	}()

	// 测试场景1：目录不存在时
	path := getDefaultConfigPath()
	if path != ymlPath {
		t.Errorf("Expected %s, got %s", ymlPath, path)
	}

	// 测试场景2：创建 .yml 文件
	os.MkdirAll(configDir, 0755)
	os.Create(ymlPath)
	path = getDefaultConfigPath()
	if path != ymlPath {
		t.Errorf("Expected %s, got %s", ymlPath, path)
	}

	// 测试场景3：创建 .yaml 文件
	os.Remove(ymlPath)
	os.Create(yamlPath)
	path = getDefaultConfigPath()
	if path != yamlPath {
		t.Errorf("Expected %s, got %s", yamlPath, path)
	}
}

func TestLoadConfig(t *testing.T) {
	// 创建临时配置文件
	content := `
- name: "test1"
  host: "localhost"
  port: 6379
  password: ""
- name: "test2"
  host: "example.com"
  port: 6380
  password: "testpass"
`
	tmpfile, err := os.CreateTemp("", "redisw_test_*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// 测试加载配置
	servers := loadConfig(tmpfile.Name())
	if len(servers) != 2 {
		t.Errorf("Expected 2 servers, got %d", len(servers))
	}

	// 验证第一个服务器配置
	if servers[0].Name != "test1" {
		t.Errorf("Expected name 'test1', got %s", servers[0].Name)
	}
	if servers[0].Port != 6379 {
		t.Errorf("Expected port 6379, got %d", servers[0].Port)
	}

	// 验证第二个服务器配置
	if servers[1].Name != "test2" {
		t.Errorf("Expected name 'test2', got %s", servers[1].Name)
	}
	if servers[1].Password != "testpass" {
		t.Errorf("Expected password 'testpass', got %s", servers[1].Password)
	}
}

func TestGetServerNames(t *testing.T) {
	servers := []RedisServer{
		{Name: "test1", Host: "localhost", Port: 6379},
		{Name: "test2", Host: "example.com", Port: 6380},
	}

	names := getServerNames(servers)
	if len(names) != 2 {
		t.Errorf("Expected 2 names, got %d", len(names))
	}
	if names[0] != "test1" {
		t.Errorf("Expected 'test1', got %s", names[0])
	}
	if names[1] != "test2" {
		t.Errorf("Expected 'test2', got %s", names[1])
	}
}