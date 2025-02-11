package lvm

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func LvmExpand(isLargeDisk bool, expansionTarget, selectedDisk string) {

	printHeader("扩容前的根分区状态：")
	runLvmExpandCommand("df", "-h", expansionTarget)

	// ext4 or xfs
	fsType, err := getFSType(expansionTarget)
	if err != nil {
		exitWithError(err)
	}

	if !isSupportedFS(fsType) {
		exitWithError(fmt.Errorf("不支持的文件系统类型: %s (仅支持 xfs/ext4)", fsType))
	}

	// lv path
	lvPath, err := getLVPath(expansionTarget)
	if err != nil {
		exitWithError(fmt.Errorf("获取逻辑卷路径失败: %v", err))
	}

	// vg name
	vgName, err := getVGName(lvPath)
	if err != nil {
		exitWithError(fmt.Errorf("获取卷组名称失败: %v", err))
	}

	// expand lvm init
	printHeader(fmt.Sprintf("开始扩容操作 [磁盘: %s 卷组: %s]", selectedDisk, vgName))
	if err := expandLVM(isLargeDisk, expansionTarget, selectedDisk, vgName, lvPath, fsType); err != nil {
		exitWithError(err)
	}

	printHeader("扩容后的根分区状态：")
	runLvmExpandCommand("df", "-h", expansionTarget)
	fmt.Println("\n✔ LVM扩容操作已成功完成")
}

func printHeader(msg string) {
	fmt.Printf("\n%s\n", msg)
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "%s %v\n", "错误:", err)
	os.Exit(1)
}

func runLvmExpandCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		exitWithError(fmt.Errorf("命令执行失败: %s %v\n原因: %v", name, arg, err))
	}
}

func getFSType(target string) (string, error) {
	cmd := exec.Command("df", "-T", target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取文件系统类型失败: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("解析df输出失败")
	}

	fields := strings.Fields(lines[len(lines)-1])
	if len(fields) < 2 {
		return "", fmt.Errorf("无法解析文件系统类型")
	}
	return fields[1], nil
}

func isSupportedFS(fsType string) bool {
	return fsType == "xfs" || fsType == "ext4"
}

func getLVPath(target string) (string, error) {
	cmd := exec.Command("df", "--output=source", target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取逻辑卷路径失败: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("无法解析逻辑卷路径")
	}
	return strings.TrimSpace(lines[1]), nil
}

func getVGName(lvPath string) (string, error) {
	cmd := exec.Command("lvs", "--noheadings", "-o", "vg_name", lvPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取卷组名称失败: %v", err)
	}

	vgName := strings.TrimSpace(string(output))
	if vgName == "" {
		return "", fmt.Errorf("卷组名称为空")
	}
	return vgName, nil
}

func expandLVM(isLargeDisk bool, expansionTarget, selectedDisk, vgName, lvPath, fsType string) error {
	newPartition := selectedDisk + "1"

	// 创建分区
	if isLargeDisk {
		fmt.Println("▷ 创建GPT分区表...")
		if err := createGPTPartition(selectedDisk); err != nil {
			return err
		}
	} else {
		fmt.Println("▷ 创建MBR分区表...")
		if err := createMBRPartition(selectedDisk); err != nil {
			return err
		}
	}

	fmt.Println("▷ 刷新分区表...")
	if err := refreshPartitions(selectedDisk); err != nil {
		return err
	}

	fmt.Println("▷ 验证新分区...")
	if err := verifyPartition(newPartition); err != nil {
		return err
	}

	fmt.Println("▷ 创建物理卷...")
	if err := createPhysicalVolume(newPartition); err != nil {
		return err
	}

	fmt.Println("▷ 扩展卷组...")
	if err := extendVolumeGroup(vgName, newPartition); err != nil {
		return err
	}

	fmt.Println("▷ 扩展逻辑卷...")
	if err := extendLogicalVolume(lvPath); err != nil {
		return err
	}

	fmt.Println("▷ 调整文件系统...")
	return resizeFilesystem(lvPath, expansionTarget, fsType)
}

func createGPTPartition(disk string) error {
	cmds := []*exec.Cmd{
		exec.Command("parted", "/dev/"+disk, "--script", "mklabel", "gpt"),
		exec.Command("parted", "/dev/"+disk, "--script", "mkpart", "primary", "0%", "100%"),
	}

	for _, cmd := range cmds {
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("分区创建失败: %v\n输出: %s", err, string(output))
		}
	}
	return nil
}

func createMBRPartition(disk string) error {
	fdiskCmd := "n\np\n1\n\n\nw\n"
	cmd := exec.Command("fdisk", "/dev/"+disk)
	cmd.Stdin = strings.NewReader(fdiskCmd)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("分区创建失败: %v\n输出: %s", err, string(output))
	}
	return nil
}

func refreshPartitions(disk string) error {
	cmd := exec.Command("partprobe", "/dev/"+disk)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("分区刷新失败: %v\n输出: %s", err, string(output))
	}
	time.Sleep(2 * time.Second) // 等待设备识别
	return nil
}

func verifyPartition(partition string) error {
	if _, err := os.Stat("/dev/" + partition); os.IsNotExist(err) {
		return fmt.Errorf("分区不存在: /dev/%s", partition)
	}
	return nil
}

func createPhysicalVolume(partition string) error {
	cmd := exec.Command("pvcreate", "/dev/"+partition)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("物理卷创建失败: %v\n输出: %s", err, string(output))
	}
	return nil
}

func extendVolumeGroup(vgName, partition string) error {
	cmd := exec.Command("vgextend", vgName, "/dev/"+partition)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("卷组扩展失败: %v\n输出: %s", err, string(output))
	}
	return nil
}

func extendLogicalVolume(lvPath string) error {
	cmd := exec.Command("lvextend", "-l", "+100%FREE", lvPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("逻辑卷扩展失败: %v\n输出: %s", err, string(output))
	}
	return nil
}

func resizeFilesystem(lvPath, mountPoint, fsType string) error {
	switch fsType {
	case "xfs":
		cmd := exec.Command("xfs_growfs", mountPoint)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("xfs扩展失败: %v\n输出: %s", err, string(output))
		}
	case "ext4":
		cmd := exec.Command("resize2fs", lvPath)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("ext4扩展失败: %v\n输出: %s", err, string(output))
		}
	}
	return nil
}
