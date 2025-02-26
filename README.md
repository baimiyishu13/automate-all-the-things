

<p title="All The Things" align="center"> <img src="./images/BbBEFEE.jpeg"> </p>



---

Automate All The Things 的最基础版本

```sh
(base) ➜  automate-all-the-things git:(master) ./padded-mac 
L 2.5

Usage:
  paddle [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  lvmCreate   lvmCreate <disk> <mountPoint> <fsType>
  lvmExpand   lvmExpand <isLargeDisk> <expansionTarget> <selectedDisk>
  rbdUseradd  rbdUseradd <user> <password> <email> <cluster> 

Flags:
  -h, --help     help for paddle
  -t, --toggle   Help message for toggle

Use "paddle [command] --help" for more information about a command.
```



---

使用：

前提: ansible

示例：扩容LVM 【lvmExpand   lvmExpand \<isLargeDisk> \<expansionTarget> \<selectedDisk>】

```sh
#!/bin/bash 
 
set -e 
 
src_script_path="/data/ml/paddle"
dest_script_path="/tmp/paddle"
 
required_vars=(isLargeDisk expansionTarget selectedDisk)
for var in "${required_vars[@]}"; do 
    [[ -z "${!var}" ]] && echo "错误: 缺少必要参数 ${var}" && exit 1 
done 
 
ansible all -i "$host_ip," -m copy -a "src=$src_script_path dest=$dest_script_path mode=0755 force=yes"
ansible all -i "$host_ip," -m shell -a "$dest_script_path lvmExpand $isLargeDisk $expansionTarget $selectedDisk"
```

---

以自动完成一些工作, 但：

请谨慎使用它！

<p title="Thanos" align="center"> <img src="./images/aa.png"> </p>
