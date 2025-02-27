package docker

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// setSysctl 函数在运行时应用内核参数。
// TODO: 考虑从配置文件中读取这些参数。
func setSysctl(param, value string) error {
	// 使用 sysctl 构建命令字符串
	cmd := fmt.Sprintf("sysctl -w %s=%s", param, value)

	// 执行 shell 命令
	if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
		return fmt.Errorf("设置 %s 失败: %v", param, err)
	}
	return nil
}

// writeDockerRepoConfig writes the Docker CE repository configuration to a YUM repository file.
func writeDockerRepoConfig() error {
	repoContent := `[docker-ce-stable]
name=Docker CE Stable - $basearch
baseurl=http://mirror.yumc.local/docker-ce/linux/centos/7/$basearch/stable
enabled=1
gpgcheck=0
`
	repoFilePath := "/etc/yum.repos.d/docker-ce.repo"

	// Write the repository content to the file
	if err := os.WriteFile(repoFilePath, []byte(repoContent), 0644); err != nil {
		return fmt.Errorf("failed to write Docker CE repository configuration: %v", err)
	}

	fmt.Println("Docker CE repository configuration written successfully")
	return nil
}

// installDocker 函数配置 sysctl，重建 RPM 数据库，安装 Docker，并更新 cgroup 设置。
// TODO: 在运行 yum install 之前检查是否已经安装了 Docker。
func InstallDocker(insecureRegistries []string, base string) error {
	// 配置网络相关的 sysctl 参数
	params := []struct {
		key   string
		value string
	}{
		//{"net.bridge.bridge-nf-call-ip6tables", "1"},
		//{"net.bridge.bridge-nf-call-iptables", "1"},
		{"net.ipv4.ip_forward", "1"},
		{"net.ipv4.conf.all.forwarding", "1"},
	}
	for _, p := range params {
		if err := setSysctl(p.key, p.value); err != nil {
			return err
		}
	}

	// 重建 RPM 数据库
	fmt.Println("==> 重建 RPM 数据库")
	if _, err := exec.Command("sh", "-c", "rpm --rebuilddb").Output(); err != nil {
		return fmt.Errorf("重建 rpm 数据库失败: %v", err)
	}

	if err := writeDockerRepoConfig(); err != nil {
		return err
	}

	// 安装 Docker CE
	fmt.Println("==> 安装 Docker CE")
	if _, err := exec.Command("sh", "-c", "yum install -y docker-ce").Output(); err != nil {
		return fmt.Errorf("安装 docker 失败: %v", err)
	}

	// 创建 daemon.json 配置文件
	fmt.Println("==> 创建 daemon.json 配置文件")
	daemonConfig := map[string]interface{}{
		"exec-opts":           []string{"native.cgroupdriver=systemd"},
		"insecure-registries": insecureRegistries,
		"storage-driver":      "overlay2",
		"default-address-pools": []map[string]interface{}{
			{
				"base": base,
				"size": 24,
			},
		},
		"log-driver": "json-file",
		"log-opts": map[string]string{
			"max-size": "100m",
			"max-file": "2",
		},
	}

	daemonConfigBytes, err := json.MarshalIndent(daemonConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("生成 daemon.json 配置失败: %v", err)
	}

	if err := os.WriteFile("/etc/docker/daemon.json", daemonConfigBytes, 0644); err != nil {
		return fmt.Errorf("写入 daemon.json 配置失败: %v", err)
	}

	// 启用并启动 Docker 服务
	fmt.Println("==> 启用并启动 Docker 服务")
	if _, err := exec.Command("sh", "-c", "systemctl daemon-reload && systemctl enable docker && systemctl start docker").Output(); err != nil {
		return fmt.Errorf("启动 docker 失败: %v", err)
	}

	//// 更新 cgroup 配置
	//fmt.Println("==> 更新 cgroup 配置")
	//cmd := "sed -i '/filebeat/d' /etc/cgrules.conf && systemctl restart cgconfig && systemctl restart cgred"
	//if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
	//	return fmt.Errorf("修改 cgroup 失败: %v", err)
	//}

	fmt.Println("==> Docker 安装完成")
	return nil
}
