package lvm

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 解析输入参数
func parseArgs(disk, mountPoint, fsType string) (string, string, string, error) {
	// 自动补全设备路径
	if !strings.HasPrefix(disk, "/dev/") {
		disk = "/dev/" + disk
	}

	if mountPoint == "/" {
		return "", "", "", errors.New("挂载点不能为根目录 (/)")
	}

	return disk, mountPoint, fsType, nil
}

// 执行系统命令
func runLvmCreateCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// 创建物理卷
func createPV(disk string) error {
	fmt.Printf("正在创建物理卷: %s...\n", disk)
	if err := runLvmCreateCommand("pvcreate", "-f", disk); err != nil {
		return fmt.Errorf("物理卷 %s 创建失败: %v", disk, err)
	}
	fmt.Printf("物理卷 %s 创建成功。\n", disk)
	return nil
}

// 创建卷组
func createVG(vgName, disk string) error {
	fmt.Printf("正在创建卷组: %s...\n", vgName)
	if err := runLvmCreateCommand("vgcreate", vgName, disk); err != nil {
		return fmt.Errorf("卷组 %s 创建失败: %v", vgName, err)
	}
	fmt.Printf("卷组 %s 创建成功。\n", vgName)
	return nil
}

// 创建逻辑卷
func createLV(vgName, lvName string) error {
	fmt.Printf("正在创建逻辑卷: %s...\n", lvName)
	if err := runLvmCreateCommand("lvcreate", "-n", lvName, "-l", "100%FREE", vgName); err != nil {
		return fmt.Errorf("逻辑卷 %s 创建失败: %v", lvName, err)
	}
	fmt.Printf("逻辑卷 %s 创建成功。\n", lvName)
	return nil
}

// 格式化逻辑卷
func formatLV(lvPath, fsType string) error {
	fmt.Printf("正在格式化逻辑卷: %s 为 %s...\n", lvPath, fsType)
	if err := runLvmCreateCommand("mkfs."+fsType, lvPath); err != nil {
		return fmt.Errorf("逻辑卷 %s 格式化失败: %v", lvPath, err)
	} else {
		fmt.Printf("逻辑卷 %s 格式化为 %s 成功。\n", lvPath, fsType)
	}
	return nil
}

// 挂载逻辑卷
func mountLV(lvPath, mountPoint string) error {
	fmt.Printf("正在挂载逻辑卷: %s 到 %s...\n", lvPath, mountPoint)
	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		return fmt.Errorf("创建挂载点 %s 失败: %v", mountPoint, err)
	}
	if err := runLvmCreateCommand("mount", lvPath, mountPoint); err != nil {
		return fmt.Errorf("挂载 %s 到 %s 失败: %v", lvPath, mountPoint, err)
	}
	fmt.Printf("逻辑卷 %s 已挂载到 %s。\n", lvPath, mountPoint)
	return nil
}

// 更新 /etc/fstab
func updateFstab(lvPath, mountPoint, fsType string) error {
	fmt.Printf("正在更新 /etc/fstab...\n")
	uuidCmd := exec.Command("blkid", "-s", "UUID", "-o", "value", lvPath)
	uuidOutput, err := uuidCmd.Output()
	if err != nil {
		return fmt.Errorf("获取 %s 的 UUID 失败: %v", lvPath, err)
	}
	uuid := strings.TrimSpace(string(uuidOutput))

	fstabEntry := fmt.Sprintf("UUID=%s  %s  %s  defaults  0  2\n", uuid, mountPoint, fsType)

	f, err := os.OpenFile("/etc/fstab", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开 /etc/fstab 失败: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(fstabEntry); err != nil {
		return fmt.Errorf("写入 /etc/fstab 失败: %v", err)
	}
	fmt.Println("/etc/fstab 已更新。")
	return nil
}

// 验证挂载
func verifyMount(mountPoint string) error {
	fmt.Printf("正在验证挂载点: %s...\n", mountPoint)
	mountCmd := exec.Command("mountpoint", "-q", mountPoint)
	if err := mountCmd.Run(); err != nil {
		return fmt.Errorf("验证失败: %s 未成功挂载", mountPoint)
	}
	fmt.Printf("验证成功: %s 已挂载。\n", mountPoint)
	return nil
}

// 设置LVM配置
func LvmCreate(disk, mountPoint, fsType string) error {
	var err error

	disk, mountPoint, fsType, err = parseArgs(disk, mountPoint, fsType)
	if err != nil {
		return err
	}

	baseName := filepath.Base(mountPoint)
	vgName := "vg" + baseName
	lvName := "lv" + baseName
	lvPath := fmt.Sprintf("/dev/%s/%s", vgName, lvName)

	fmt.Println("开始设置LVM配置...")
	fmt.Printf("磁盘设备: %s\n挂载点: %s\n文件系统类型: %s\n", disk, mountPoint, fsType)

	if err = createPV(disk); err != nil {
		return fmt.Errorf("创建物理卷失败: %v", err)
	}

	if err = createVG(vgName, disk); err != nil {
		return fmt.Errorf("创建卷组失败: %v", err)
	}

	if err = createLV(vgName, lvName); err != nil {
		return fmt.Errorf("创建逻辑卷失败: %v err")
	}

	if err = formatLV(lvPath, fsType); err != nil {
		return fmt.Errorf("格式化逻辑卷失败: %v", err)
	}

	if err = mountLV(lvPath, mountPoint); err != nil {
		return fmt.Errorf("挂载逻辑卷失败: %v", err)
	}

	if err = updateFstab(lvPath, mountPoint, fsType); err != nil {
		return fmt.Errorf("更新fstab失败: %v", err)
	}

	if err = verifyMount(mountPoint); err != nil {
		return fmt.Errorf("验证挂载失败: %v", err)
	}

	fmt.Println("LVM配置设置完成！")
	return nil
}
