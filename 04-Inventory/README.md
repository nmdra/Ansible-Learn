# Inventory

> ‘Chapter 4. Inventory: Describing Your Servers’
—Bas Meijer, Lorin Hochstein, Rene Moser, “Ansible: Up and Running”

The collection of hosts that Ansible knows about is called the _inventory_.

## Behavioral inventory parameters

In Ansible, **behavioral inventory parameters** are configuration options that determine how Ansible interacts with hosts in the inventory. These parameters control aspects such as connection methods, user privileges, and specific runtime settings.

| Name                           | Default           | Description                                                                       |
| ------------------------------ | ----------------- | --------------------------------------------------------------------------------- |
| `ansible_host`                 | Name of host      | Hostname or IP address to SSH to                                                  |
| `ansible_port`                 | 22                | Port to SSH to                                                                    |
| `ansible_user`                 | $USER             | User to SSH as                                                                    |
| `ansible_password`             | _(None)_          | Password to use for SSH authentication                                            |
| `ansible_connection`           | smart             | How Ansible will connect to host (see the following section)                      |
| `ansible_ssh_private_key_file` | _(None)_          | SSH private key to use for SSH authentication                                     |
| `ansible_shell_type`           | sh                | Shell to use for commands (see the following section)                             |
| `ansible_python_interpreter`   | _/usr/bin/python_ | Python interpreter on host (see the following section)                            |
| `ansible_*_interpreter`        | (_None_)          | Like `ansible_python_interpreter` for other languages (see the following section) |


You can override some of the behavioral parameter default values in the inventory file, or you can override them in the `defaults` section of the _ansible.cfg_ file.

## Groups

### Basic Groups and Group Vars

```ini
[webservers]
web1 ansible_host=192.168.1.10
web2 ansible_host=192.168.1.11

[databases]
db1 ansible_host=192.168.1.20
db2 ansible_host=192.168.1.21

[webservers:vars]
ansible_user=ubuntu
ansible_become=true

[databases:vars]
ansible_user=postgres
ansible_become=true
ansible_python_interpreter=/usr/bin/python3
```

### Nested Groups

```ini
[webservers]
web1
web2

[databases]
db1
db2

[production:children]
webservers
databases

```

### Numbered Hosts

```ini
[web] 
web[1:20].example.com
```

```ini
[web] 
web-[a:t].example.com
```

## Dynamic Inventory

Dynamic Inventory is a feature in Ansible that enables automatic fetching of host information from external sources, such as cloud platforms, containers, or custom APIs.
(This approach eliminates the need to maintain a static inventory file, making it ideal for environments where hosts are frequently added, removed, or modified.)

**If the inventory file is marked executable, Ansible will assume it is a dynamic inventory script and will execute the file instead of reading it.**

## Inventory Plugins

Ansible inventory plugins are components that allow you to dynamically define and manage inventories in various formats and sources.

1. **Built-in Inventory Plugins**: Ansible comes with several built-in inventory plugins that support common use cases, such as fetching inventory data from files, cloud providers, or APIs.
2. **Custom Inventory Plugins**: Users can develop custom plugins to meet specific requirements or integrate with proprietary systems.

To see the list of available plug-ins:

```bash
$ ansible-doc -t inventory -l
```

### Common Built-in Inventory Plugins

| Plugin Name       | Description                                                                                  |
| ----------------- | -------------------------------------------------------------------------------------------- |
| **`ini`**         | Parses inventory from `ini`-formatted files (default inventory format).                      |
| **`yaml`**        | Parses inventory from `YAML` files.                                                          |
| **`host_list`**   | Reads inventory from a list of hosts provided inline.                                        |
| **`script`**      | Executes a custom script that generates inventory dynamically.                               |
| **`directory`**   | Loads inventory from multiple files in a directory.                                          |
| **`cloud`**       | Retrieves inventory from cloud providers (e.g., AWS, GCP, Azure, etc.).                      |
| **`constructed`** | Creates a dynamic inventory by filtering or transforming data from another inventory source. |
| **`openstack`**   | Fetches inventory from OpenStack.                                                            |
| **`azure_rm`**    | Retrieves inventory from Azure Resource Manager.                                             |
| **`aws_ec2`**     | Gathers inventory from AWS EC2 instances.                                                    |
| **`gcp_compute`** | Fetches inventory from Google Cloud Platform Compute Engine.                                 |
| **`vmware_vm`**   | Retrieves inventory from VMware environments.                                                |
### Docker Dynamic Inventory

```bash
ansible-galaxy collection install community.docker
```

```yaml
plugin: community.docker.docker
compose: yes
containers:
  - name: app_container
```

```bash
ansible-inventory -i docker_inventory.yml --list
```

## Directory-Based Group Variables

Ansible handles **inventory management** when using a combination of **static inventory files** and **dynamic inventory scripts**.

`ansible.cfg`
```ini
[defaults]
inventory = ./inventory
```

```bash
ansible-playbook -i ./inventory playbook.yml
```

```plaintext
inventory/
├── aws_ec2.yml        # Dynamic inventory for AWS EC2
├── azure_rm.yml       # Dynamic inventory for Azure Resource Manager
├── group_vars/        # Variables specific to groups
│   ├── vagrant        # Variables for Vagrant hosts
│   ├── staging        # Variables for staging environment
│   ├── production     # Variables for production environment
├── host_vars/         # Variables specific to individual hosts
│   ├── host1          # Variables for "host1"
│   ├── host2          # Variables for "host2"
├── hosts              # Static inventory for generic hosts
├── vagrant.py         # Dynamic inventory script for Vagrant
```

```yaml
# inventory/host_vars/webserver1
ansible_host: 192.168.1.100
ansible_user: ubuntu
app_environment: production
app_port: 8080
```

## Adding Entries at Runtime

### add_host

The `add_host` module adds a host to the inventory; this is useful if you’re using Ansible to provision new virtual machine instances inside an infrastructure-as-a-service cloud.

```yaml
---
- name: Dynamically add a host to inventory
  hosts: localhost
  tasks:
    - name: Provision a new VM (simulated)
      command: echo "192.168.56.101"
      register: new_vm_ip

    - name: Add the new VM to the inventory
      ansible.builtin.add_host:
        name: "{{ new_vm_ip.stdout }}"
        groups: new_vms
        ansible_user: vagrant
        ansible_port: 2222
	ansible_private_key_file: >
		.vagrant/machines/default/virtualbox/private_key

- name: Configure the newly added VM
  hosts: new_vms
  tasks:
    - name: Ping the new VM
      ansible.builtin.ping:
```
### group_by

Used to create groups based on host variables during runtime.
```yaml
---
- name: Group hosts by their OS distribution
  hosts: all
  tasks:
    - name: Collect OS facts
      ansible.builtin.setup:

    - name: Group hosts by distribution
      ansible.builtin.group_by:
        key: "os_{{ ansible_facts['distribution'] | lower }}"

- name: Perform tasks on Ubuntu hosts
  hosts: os_ubuntu
  tasks:
    - name: Install Nginx
      ansible.builtin.apt:
        name: nginx
        state: present
        update_cache: yes
      when: ansible_facts['distribution'] == 'Ubuntu'

- name: Perform tasks on CentOS hosts
  hosts: os_centos
  tasks:
    - name: Install Nginx
      ansible.builtin.yum:
        name: nginx
        state: present
```

