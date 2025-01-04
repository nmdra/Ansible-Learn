# Vagrant Setup using Rocky Linux 9.5

This Vagrant setup is based on Rocky Linux 9.5  for Ansible testing.

> [!IMPORTANT]  
> **Resource Links**: [bookmarks.nimendra.xyz/vagrant](https://bookmarks.nimendra.xyz/bookmarks/shared?q=%23vagrant)  

---

## **Guide**

> [!WARNING]
> Make sure to install [Vagrant](https://www.vagrantup.com/downloads) and [VirtualBox](https://www.virtualbox.org/wiki/Downloads) before starting. 

### **Download Rocky Linux Vagrant Box**

To explore other versions, visit: [Rocky Linux 9.5 Vault](https://dl.rockylinux.org/vault/rocky/9.5/images/x86_64/).


### **Steps to Set Up**

1. **Create a Project Directory**  
   ```bash
   mkdir vagrant && cd vagrant
   ```

2. **Download the Rocky Linux Box**  
   ```bash
   wget https://dl.rockylinux.org/pub/rocky/9/images/x86_64/Rocky-9-Vagrant-Vbox.latest.x86_64.box
   ```

3. **Add the Box to Vagrant**  
   ```bash
   vagrant box add rockybox Rocky-9-Vagrant-Vbox.latest.x86_64.box
   ```

4. **Initialize the Vagrantfile**  
   ```bash
   vagrant init rockybox
   ```

5. **Start the Vagrant Machine**  
   ```bash
   vagrant up
    ```

---

## **Common Commands**

- **SSH into the VM**:  
  ```bash
  vagrant ssh
  ```

- **Restart the VM (apply changes)**:  
  ```bash
  vagrant reload
  ```

- **Stop the VM**:  
  ```bash
  vagrant halt
  ```

- **Destroy the VM**:  
  ```bash
  vagrant destroy
  ```
