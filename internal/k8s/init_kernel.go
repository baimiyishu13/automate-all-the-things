package k8s

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"
)

// 优化内核参数
func optimizeLinux() error {
	fmt.Println("==> 正在优化内核参数...")

	// 写入 /etc/sysctl.conf
	data := []byte(`
net.bridge.bridge-nf-call-ip6tables=1
net.bridge.bridge-nf-call-iptables=1
net.ipv4.ip_forward=1
net.ipv4.conf.all.forwarding=1
net.ipv4.neigh.default.gc_thresh1=4096
net.ipv4.neigh.default.gc_thresh2=6144
net.ipv4.neigh.default.gc_thresh3=8192
net.ipv4.neigh.default.gc_interval=60
net.ipv4.neigh.default.gc_stale_time=120

# 参考 https://github.com/prometheus/node_exporter#disabled-by-default
kernel.perf_event_paranoid=-1
kernel.pty.max=8192

#sysctls for k8s node config
net.ipv4.tcp_slow_start_after_idle=0
net.core.rmem_max=16777216
fs.inotify.max_user_watches=524288
kernel.softlockup_all_cpu_backtrace=1

kernel.softlockup_panic=0

kernel.watchdog_thresh=30
fs.file-max=50607016
fs.inotify.max_user_instances=81920
fs.inotify.max_queued_events=16384
vm.max_map_count=262144
fs.may_detach_mounts=1
net.core.netdev_max_backlog=16384
net.ipv4.tcp_wmem=4096 12582912 16777216
net.core.wmem_max=16777216
net.core.somaxconn=32768
net.ipv4.ip_forward=1
net.ipv4.tcp_max_syn_backlog=8096
net.ipv4.tcp_rmem=4096 12582912 16777216

net.ipv6.conf.all.disable_ipv6=1
net.ipv6.conf.default.disable_ipv6=1
net.ipv6.conf.lo.disable_ipv6=1

kernel.yama.ptrace_scope=0
vm.swappiness=0

# 可以控制core文件的文件名中是否添加pid作为扩展。
kernel.core_uses_pid=1

# Do not accept source routing
net.ipv4.conf.default.accept_source_route=0
net.ipv4.conf.all.accept_source_route=0

# Promote secondary addresses when the primary address is removed
net.ipv4.conf.default.promote_secondaries=1
net.ipv4.conf.all.promote_secondaries=1

# Enable hard and soft link protection
fs.protected_hardlinks=1
fs.protected_symlinks=1

# 源路由验证
# see details in https://help.aliyun.com/knowledge_detail/39428.html
net.ipv4.conf.all.rp_filter=0
net.ipv4.conf.default.rp_filter=0
net.ipv4.conf.default.arp_announce = 2
net.ipv4.conf.lo.arp_announce=2
net.ipv4.conf.all.arp_announce=2

# see details in https://help.aliyun.com/knowledge_detail/41334.html
net.ipv4.tcp_max_tw_buckets=5000
net.ipv4.tcp_syncookies=1
net.ipv4.tcp_fin_timeout=30
net.ipv4.tcp_synack_retries=2
kernel.sysrq=1
`)
	if err := ioutil.WriteFile(`/etc/sysctl.conf`, data, 0644); err != nil {
		return fmt.Errorf("写入 sysctl.conf 失败: %v", err)
	}

	// 执行 sysctl -p
	cmd := exec.Command("sh", "-c", "sysctl -p >/dev/null 2>&1")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("执行 sysctl -p 失败: %v", err)
	}

	fmt.Println("==> 已完成内核参数优化")
	return nil
}

// 优化文件描述符限制
func optimizeLimits() error {
	fmt.Println("==> 正在优化文件描述符限制...")

	// 将内容追加到 /etc/security/limits.conf
	limitsData := []byte(`
* soft nofile 1024000
* hard nofile 1024000
`)
	f, err := os.OpenFile(`/etc/security/limits.conf`, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开 limits.conf 失败: %v", err)
	}
	defer f.Close()

	if _, err = f.Write(limitsData); err != nil {
		return fmt.Errorf("写入 limits.conf 失败: %v", err)
	}

	fmt.Println("==> 已完成文件描述符限制优化")
	return nil
}

func chmod(path string, mode uint32) error {
	return syscall.Chmod(path, mode)
}

func KernelInit() {
	// 临时将权限调高
	if err := chmod(`/etc/sysctl.conf`, 0777); err != nil {
		log.Fatalf("修改 /etc/sysctl.conf 权限失败: %v", err)
	}
	if err := chmod(`/sbin/sysctl`, 0777); err != nil {
		log.Fatalf("修改 /sbin/sysctl 权限失败: %v", err)
	}
	if err := chmod(`/etc/security/limits.conf`, 0777); err != nil {
		log.Fatalf("修改 /etc/security/limits.conf 权限失败: %v", err)
	}

	// 执行优化
	if err := optimizeLinux(); err != nil {
		log.Fatal(err)
	}
	if err := optimizeLimits(); err != nil {
		log.Fatal(err)
	}

	// 恢复默认权限
	if err := chmod(`/etc/sysctl.conf`, 0644); err != nil {
		log.Fatalf("还原 /etc/sysctl.conf 权限失败: %v", err)
	}
	if err := chmod(`/sbin/sysctl`, 0755); err != nil {
		log.Fatalf("还原 /sbin/sysctl 权限失败: %v", err)
	}
	if err := chmod(`/etc/security/limits.conf`, 0644); err != nil {
		log.Fatalf("还原 /etc/security/limits.conf 权限失败: %v", err)
	}

	fmt.Println("系统优化已完成。")
}
