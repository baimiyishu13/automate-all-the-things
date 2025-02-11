package rbd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateUserRequest 结构体定义
type CreateUserRequest struct {
	UserName     string `json:"nick_name"`
	Password     string `json:"password"`
	EnterpriseId string `json:"enterprise_id"`
	Email        string `json:"email"`
}

// GetApiURL 根据环境和API功能生成完整的API URL
func GetApiURL(environment, apiFunction string) (string, error) {
	baseURL, ok := EnvApiURLMap[environment]
	if !ok {
		return "", fmt.Errorf("未找到对应的环境：%s", environment)
	}

	path, ok := ApiPathMap[apiFunction]
	if !ok {
		return "", fmt.Errorf("未找到对应的API功能：%s", apiFunction)
	}

	return fmt.Sprintf("%s%s", baseURL, path), nil
}

// CreateUser 函数用于创建用户
func CreateUser(username, password, email, cluster string) error {
	// 获取API URL
	apiURL, err := GetApiURL(cluster, "user_create")
	//fmt.Println(apiURL)
	if err != nil {
		return fmt.Errorf("获取API URL失败：%v", err)
	}

	// 获取集群ID和Token
	enterpriseId, ok := EnvClusterMap[cluster]
	if !ok {
		return fmt.Errorf("未找到对应的集群ID，请检查输入的cluster是否正确")
	}

	//fmt.Println("enterpriseId:", enterpriseId)

	token, ok := EnvTokenMap[cluster]
	if !ok {
		return fmt.Errorf("未找到对应的Token，请检查输入的cluster是否正确")
	}

	//fmt.Println("token:", token)

	userRequest := CreateUserRequest{
		UserName:     username,
		Password:     password,
		EnterpriseId: enterpriseId,
		Email:        email,
	}

	requestBody, err := json.Marshal(userRequest)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	//fmt.Println("json:", bytes.NewBuffer(requestBody))

	// 创建HTTP请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		panic("创建请求失败: " + err.Error())
	}

	//fmt.Println(bytes.NewBuffer(requestBody))

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic("请求发送失败: " + err.Error())
	}
	defer resp.Body.Close()

	// 处理响应
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("用户 %s 创建成功！\n", username)
	} else {
		fmt.Printf("用户创建失败，状态码：%d\n", resp.StatusCode)

		var errorBody bytes.Buffer
		errorBody.ReadFrom(resp.Body)
		fmt.Println("错误响应：", errorBody.String())
	}
	return nil
}
