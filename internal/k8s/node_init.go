package k8s

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// 检查数据盘挂载
func checkAndMountDataDisk() error {
	fmt.Println("==> 1. 检查数据盘挂载")
	// 检查根是否挂载
	cmd := "mount | grep ' on / '"
	if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
		return fmt.Errorf("根未挂载: %v", err)
	}

	// 检查根的可用空间是否大于5TB
	cmd = "df -h / | awk 'NR==2 {print $4}'"
	output, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return fmt.Errorf("检查根可用空间失败: %v", err)
	}

	availableSpace := strings.TrimSpace(string(output))
	if !strings.HasSuffix(availableSpace, "T") || strings.TrimSuffix(availableSpace, "T") < "5" {
		return fmt.Errorf("根可用空间不足5TB: %s", availableSpace)
	}
	return nil
}

// 修改节点内核参数
func modifyKernelParameters() error {
	fmt.Println("==> 2. 修改节点内核参数")
	// 这里假设需要修改的内核参数在/etc/sysctl.conf中
	cmd := "echo 'net.ipv4.ip_forward = 1' >> /etc/sysctl.conf && sysctl -p"
	if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
		return fmt.Errorf("修改内核参数失败: %v", err)
	}
	return nil
}

// 检查DNS配置
func checkDNSConfig() error {
	fmt.Println("==> 3. 检查DNS配置")
	// 检查/etc/resolv.conf是否包含正确的DNS配置
	cmd := "grep 'nameserver' /etc/resolv.conf"
	output, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil || !strings.Contains(string(output), "8.8.8.8") {
		return fmt.Errorf("DNS配置错误")
	}
	return nil
}

// 设置主机名
func setHostname(hostname string) error {
	fmt.Println("==> 4. 设置主机名")
	cmd := fmt.Sprintf("hostnamectl set-hostname --static %s", hostname)
	if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
		return fmt.Errorf("设置主机名失败: %v", err)
	}
	return nil
}

// 网卡摘除DNS相关信息
func removeDNSFromNetworkConfig() error {
	fmt.Println("==> 5. 网卡摘除DNS相关信息")
	// 这里假设网卡配置文件是/etc/sysconfig/network-scripts/ifcfg-bond0
	cmd := "sed -i '/^DNS/d' /etc/sysconfig/network-scripts/ifcfg-bond0"
	if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
		return fmt.Errorf("移除DNS信息失败: %v", err)
	}
	return nil
}

// 执行安装Docker的Jenkins任务
func runJenkinsJob() error {
	fmt.Println("==> 6. 执行安装Docker的Jenkins任务")
	// 这里假设Jenkins任务的URL是http://jenkins/job/install-docker/build
	cmd := "curl -X POST http://jenkins/job/install-docker/build"
	if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
		return fmt.Errorf("执行Jenkins任务失败: %v", err)
	}
	return nil
}

// 检查网络连通性
func checkNetworkConnectivity() error {
	fmt.Println("==> 7. 检查网络连通性")
	// 这里假设需要ping的目标是同机架、同机房、跨单元、跨机房、外网主机
	targets := []string{"同机架IP", "同机房IP", "跨单元IP", "跨机房IP", "外网IP"}
	for _, target := range targets {
		cmd := fmt.Sprintf("ping -c 3 %s", target)
		if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
			return fmt.Errorf("无法ping通目标: %s", target)
		}
	}
	return nil
}

// 检查主机配置
func checkHostConfig() error {
	fmt.Println("==> 9. 检查主机配置")
	// 关闭Selinux
	cmd := "setenforce 0"
	if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
		return fmt.Errorf("关闭Selinux失败: %v", err)
	}

	// 关闭防火墙
	cmd = "systemctl stop firewalld && systemctl disable firewalld"
	if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
		return fmt.Errorf("关闭防火墙失败: %v", err)
	}

	// 关闭swap
	cmd = "swapoff -a"
	if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
		return fmt.Errorf("关闭swap失败: %v", err)
	}

	// 配置时间服务器
	cmd = "echo 'server ntp.example.com' >> /etc/chrony.conf && systemctl restart chronyd"
	if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
		return fmt.Errorf("配置时间服务器失败: %v", err)
	}

	return nil
}

// 初始化节点
func InitNode(hostname string) {
	if err := checkAndMountDataDisk(); err != nil {
		log.Fatalf("数据盘挂载检查失败: %v", err)
	}

	if err := modifyKernelParameters(); err != nil {
		log.Fatalf("节点内核参数修改失败: %v", err)
	}

	if err := checkDNSConfig(); err != nil {
		log.Fatalf("DNS配置检查失败: %v", err)
	}

	if err := setHostname(hostname); err != nil {
		log.Fatalf("设置主机名失败: %v", err)
	}

	if err := removeDNSFromNetworkConfig(); err != nil {
		log.Fatalf("网卡摘除DNS信息失败: %v", err)
	}

	if err := runJenkinsJob(); err != nil {
		log.Fatalf("安装Docker失败: %v", err)
	}

	if err := checkNetworkConnectivity(); err != nil {
		log.Fatalf("网络连通性检查失败: %v", err)
	}

	if err := checkHostConfig(); err != nil {
		log.Fatalf("主机配置检查失败: %v", err)
	}

	fmt.Println("节点初始化完成")
}
