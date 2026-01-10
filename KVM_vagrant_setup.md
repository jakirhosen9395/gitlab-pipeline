# 🚀 Create Virtual Machines with Vagrant + KVM/libvirt on Ubuntu

> **A clean, eye‑catchy, production‑ready guide** to running local VMs using **Vagrant** backed by **KVM/libvirt**.  
> This document is based on real hands‑on steps executed on a local Ubuntu system, structured to avoid confusion and guesswork.

---

## 🧭 Architecture at a Glance

```
Vagrant CLI
   │
   ▼
vagrant-libvirt provider
   │
   ▼
libvirt (libvirtd)
   │
   ▼
KVM / QEMU  →  Virtual Machine (Ubuntu 24.04)
```

---

## 🎯 What You’ll Achieve

By the end of this guide, you will have:

- ✅ Hardware‑accelerated virtualization (KVM) enabled
- ✅ libvirt configured and accessible without root
- ✅ Vagrant installed from the official HashiCorp repo
- ✅ `vagrant-libvirt` provider working
- ✅ A reproducible **Ubuntu 24.04 VM**
- ✅ Private networking with a static IP

---

## 📚 Useful Resources (bookmark these)

- 🔗 Vagrant official site:  
  https://www.vagrantup.com/

- 🔗 Vagrant box discovery (official):  
  https://portal.cloud.hashicorp.com/vagrant/discover/

- 🔗 Bento Ubuntu boxes:  
  https://portal.cloud.hashicorp.com/vagrant/discover/bento

- 🔗 libvirt docs:  
  https://libvirt.org/

- 🔗 vagrant-libvirt plugin:  
  https://github.com/vagrant-libvirt/vagrant-libvirt

---

## 🧱 Prerequisites

| Requirement | Notes |
|------------|------|
| OS | Ubuntu **24.04 LTS** (recommended) |
| Access | `sudo` privileges |
| CPU | Intel VT‑x or AMD‑V |
| BIOS | Virtualization **enabled** |

---

## 1️⃣ System Preparation (Clean Baseline)

Update the system to a known‑good state:

```bash
sudo apt update -y && sudo apt upgrade -y && sudo apt dist-upgrade -y && sudo apt autoremove -y && sudo apt clean
```

### ✅ Example output
```text
0 upgraded, 0 newly installed, 0 to remove and 0 not upgraded.
```

---

## 2️⃣ Reboot Safety Check

```bash
test -f /var/run/reboot-required && echo "🚨 Reboot required" || echo "✅ No reboot required"
```

If required:
```bash
sudo reboot now
```

---

## 3️⃣ Verify User Context

```bash
whoami
```

### ✅ Example
```text
jakir
```

---

## 4️⃣ Install & Enable SSH (Recommended)

SSH makes VM and host access predictable and script‑friendly.

```bash
sudo apt install -y openssh-server
sudo systemctl status ssh
sudo systemctl enable ssh
```

### ✅ Expected state
```text
Active: active (running)
```

---

## 5️⃣ Host Capability Checks (CPU + Virtualization)

```bash
timedatectl
lscpu
```

Check virtualization flags:

```bash
egrep -c '(vmx|svm)' /proc/cpuinfo
```

### ✅ Example
```text
8
```

Confirm OS sees virtualization:

```bash
lscpu | grep Virtualization
```

### ✅ Example
```text
Virtualization: VT-x
```

---

## 6️⃣ Install KVM + libvirt Stack

```bash
sudo apt install -y   qemu-kvm   libvirt-daemon-system   libvirt-clients   virt-manager   bridge-utils
```

Enable and start libvirt:

```bash
sudo systemctl enable libvirtd
sudo systemctl start libvirtd
systemctl status libvirtd
```

### ✅ Expected
```text
Active: active (running)
```

Reboot if prompted.

---

## 7️⃣ Fix Permissions (Critical Step ⚠️)

Without this, Vagrant **will fail**.

```bash
sudo usermod -aG libvirt $USER
sudo usermod -aG libvirt-kvm $USER
sudo usermod -aG libvirt,kvm $USER
```

Apply changes:
```bash
newgrp libvirt
newgrp kvm
```

Validate:
```bash
groups
ls -l /dev/kvm
virsh list --all
```

### ✅ `/dev/kvm` should look like
```text
crw-rw----+ 1 root kvm ...
```

---

## 8️⃣ Install Vagrant (Official HashiCorp Repo)

```bash
sudo apt install -y curl gnupg software-properties-common
```

Add HashiCorp GPG key:
```bash
curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
```

Add repo and install:
```bash
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list

sudo apt update
sudo apt install -y vagrant
```

Verify:
```bash
vagrant version
```

---

## 9️⃣ Install `vagrant-libvirt` Provider

```bash
sudo apt install -y libvirt-dev ruby-dev make gcc
vagrant plugin install vagrant-libvirt
vagrant plugin list
```

### ✅ Expected
```text
vagrant-libvirt (x.x.x, global)
```

---

## 🔟 Create Project & Initialize Box

```bash
mkdir -p ~/vagrant/ubuntu-24.04
cd ~/vagrant/ubuntu-24.04
vagrant init bento/ubuntu-24.04
```

---

## 1️⃣1️⃣ Recommended `Vagrantfile`

```ruby
Vagrant.configure("2") do |config|
  vm_ip = "192.168.169.100"

  config.vm.box = "bento/ubuntu-24.04"
  config.vm.hostname = "ubuntu24-vagrant"
  config.vm.network "private_network", ip: vm_ip

  config.vm.provider :libvirt do |libvirt|
    libvirt.driver = "kvm"
    libvirt.memory = 2048
    libvirt.cpus = 2
    libvirt.nested = true
  end
end
```

---

## 1️⃣2️⃣ VM Lifecycle Commands

### ▶️ Start
```bash
vagrant up --provider=libvirt
```

### 📊 Status
```bash
vagrant status
```

### 🔐 SSH
```bash
vagrant ssh
```

### 🔄 Reload
```bash
vagrant reload
```

### 💥 Destroy
```bash
vagrant destroy -f
```

---

## 1️⃣3️⃣ Global Cleanup & Control

```bash
vagrant global-status
vagrant destroy <id>
rm -rf .vagrant
```

---

## 🏁 Final Check

If this works:

```bash
vagrant status
vagrant ssh
```

🎉 **Congrats — your local KVM‑powered VM stack is production‑grade and reproducible.**
