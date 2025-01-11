# Variables

## Define Variables

```yaml
vars:
	tls_dir: /etc/nginx/ssl/
	key_file: nginx.key
	server_name: localhost
```

### Using var file

```yaml
vars_files:
	- nginx.yml
```

```yaml
# nginx.yml
tls_dir: /etc/nginx/ssl/
key_file: nginx.key
server_name: localhost
```

## Variable Interpolation

```yaml
- name: Varible
  debug:
	  msg: "The file used was {{ conf_file }}"
```

```yaml
- name: Concatanate variable
  debug:
	  msg: "The URL is https://{{ sever_name ~'.'~ domain +_name }}"
```

## Registering Variables

Ansible module returns results in JSON format. To use these results, you create a _registered variable_ using the `register` clause when invoking a module.

```yaml
- name: Capture and display the current user
  hosts: localhost
  tasks:
    - name: Capture output of whoami
      ansible.builtin.command: whoami
      register: login

    - name: Show output
      ansible.builtin.debug:
        msg: "The user running this playbook is: {{ login.stdout }}"
```

## Accessing Dictionary Keys

If a variable contains a dictionary, you can access the keys of the dictionary by using either a dot (`.`) or a subscript (`[]`).
```yaml
{{ result.stat }}
```

We could have used subscript notation instead:
```yaml
{{ result['stat'] }}
```

This rule applies to multiple de-references, so all of the following are equivalent:

```yaml
result['stat']['mode']
result['stat'].mode
result.stat['mode']
result.stat.mode
```
# Facts

When Ansible gathers facts, it connects to the hosts and queries it for all kinds of details about the hosts: CPU architecture, operating system, IP addresses, memory info, disk info, and more.

```yaml
- name: 'Ansible facts.'
  hosts: all
  gather_facts: true
  tasks:
	  - name: Print out OS details
  debug:
	  msg: >-
		  os_family: {{ ansible_facts.os_family }},
		  distro: 
		  {{ ansible_facts.distribution }} 
		  {{ ansible_facts.distribution_version }}, 
		  kernel: 
		  {{ ansible_facts.kernel }
```

```bash
# View all facts
ansible server -m setup

# Viewing a Subset of Facts
ansible all -m setup -a 'filter=ansible_all_ipv6_addresses'
```

## Module Generated Facts

- Some Ansible modules return a dictionary with the key `ansible_facts`.
- These facts are automatically associated with the active host as variables.
- Modules ending with `_info` return information about non-unique host objects (e.g., services, packages).

```yaml
- name: Get services facts
  service_facts:

- name: Debug SSHD service state
  debug: var=ansible_facts['services']['sshd.service']
```

## Local Facts

- Files placed in `/etc/ansible/facts.d` on a host are treated as custom facts.
- Supported formats:
    - `.ini`
    - JSON
    - Executables outputting JSON

```yaml
- name: Print book title
  debug: msg="The title is {{ ansible_local.example.book.title }}"
```

## Set Facts

- Defines new variables dynamically within a playbook.
- Useful for simplifying variable references or defining new facts based on conditions.

```yaml
- name: Set nginx state
  when: ansible_facts.services.nginx.state is defined
  set_fact:
    nginx_state: "{{ ansible_facts.services.nginx.state }}"
```

# Magic / Built-in variables

🌟 [Discovering variables: facts and magic variables — Ansible Community Documentation](https://docs.ansible.com/ansible/latest/playbook_guide/playbooks_vars_facts.html#information-about-ansible-magic-variables)

| Parameter                  | Description                                                                                                                                                                                           |
| -------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `hostvars`                 | A dict whose keys are Ansible hostnames and values are dicts that map variable names to values                                                                                                        |
| `inventory_hostname`       | The name of the current host as known in the Ansible inventory, might include domain name                                                                                                             |
| `inventory_hostname_short` | Name of the current host as known by Ansible, without the domain name (e.g., `myhost`)                                                                                                                |
| `group_names`              | A list of all groups that the current host is a member of                                                                                                                                             |
| `groups`                   | A dict whose keys are Ansible group names and values are a list of hostnames that are members of the group. Includes `all` and ungrouped groups: {“`all`”: [...], “web”: [...], “`ungrouped`”: [...]} |
| `ansible_check_mode`       | A boolean that is true when running in check mode (see “Check Mode”)                                                                                                                                  |
| `ansible_play_batch`       | A list of the inventory hostnames that are active in the current batch (see “Running on a Batch of Hosts at a Time”)                                                                                  |
| `ansible_play_hosts`       | A list of all of the inventory hostnames that are active in the current play                                                                                                                          |
| `ansible_version`          | A dict with Ansible version info: {“`full`”: `2.3.1.0`”, “`major`”: `2`, “`minor`”: `3`, “`revision`”: `1`, “`string`”: “`2.3.1.0`”}                                                                  |

- `hostvars` is computed when you run Ansible, while `host_vars` is a directory that you can use to define variables for a particular system.
# Variable Precedence

🌟 [Using Variables](https://docs.ansible.com/ansible/latest/playbook_guide/playbooks_variables.html#understanding-variable-precedence)

---
## 1. The `when` Statement

- **Purpose**: Executes tasks conditionally based on variables or previous task outputs.
- **Examples**:
    - Use boolean variables (`when: is_db_server`).
    - Handle undefined variables (`when: is_db_server is defined and is_db_server`).
    - Evaluate registered variables (`when: 'ready' in myapp_result.stdout`).

## 2. `changed_when` and `failed_when`

- **`changed_when`**: Overrides the default behavior for determining task changes.
    - Example: Detect changes based on specific output (e.g., `'Nothing to install' not in composer.stdout`).
- **`failed_when`**: Customizes failure conditions.
    - Example: Ignore errors if a specific string is present in `stderr`.

## 3. `ignore_errors`

- **Purpose**: Prevents Ansible from failing when a task encounters errors.
- **Caution**: Use sparingly, as it may hide actual issues.

## 4. Delegation and Local Actions

- **Delegating Tasks**: Use `delegate_to` to run tasks on specific hosts.
    - Example: Send commands to a monitoring server or load balancer.
- **Local Execution**: Use `local_action` as shorthand for tasks run on the control machine.
    - Example: Manage load balancers or check system states locally.

## 5. Pausing Execution with `wait_for`

- **Purpose**: Wait for resources (e.g., ports, files, or drained connections) to become available.
- **Features**:
    - Wait for ports to open or close.
    - Check for file presence or absence.
    - Introduce delays for synchronization.

## 6. Running Playbooks Locally

- **Use Case**: Run playbooks on the same machine as the Ansible control host using `--connection=local`.
- **Benefits**:
    - Speeds up execution by avoiding SSH overhead.
    - Useful for self-provisioning or CI/CD pipelines.
- **Example**:

```yaml
- hosts: 127.0.0.1
  gather_facts: no
  tasks:
    - name: Get the current date
      command: date
      register: date
    - debug:
        var: date.stdout
```

# Extras

## Prompts

- `private`: If set to `yes`, the user’s input will be hidden on the command line.
- `default`: You can set a default value for the prompt, to save time for the end user.

```yaml
- hosts: all

vars_prompt:
	- name: Username
	  prompt: "Enter Your Username?"
	  private: False
```

## Tags

Tags allow you to run (or exclude) subsets of a playbook’s tasks.
You can tag roles, included files, individual tasks, and even entire plays.

```yaml
- hosts: servers
  tags: deploy

  roles:
	- role: tomcat
	    tags: ['tomcat','app']
  tasks:
	- name: Notify
	  local_action: 
		  module: osx_say
		  msg: "{{inventory_hostname}} is finished"
	  tags:
		  - notification
```

adding more than one tag, you have to use YAML’s list syntax, for example:

```yaml
# shorthand list
tags: ['red', 'green']

# Explicit List
tags:
	- red
	- green
```

## Blocks

**blocks** are a way to group a set of tasks together in a playbook. Blocks allow you to apply conditional checks, error handling, and other features to a group of tasks at once, making your playbooks more organized and efficient.

```yaml
- hosts: web
  tasks:
    - block:
	    - dnf: name=https state=present
	    - template: src=httpd.conf.j2 dest=/etc/httpd/conf/httpd.conf
	    - service: name=httpd state=started enabled=yes
	  when: ansible_os_family = 'RedHat'
	  become: yes

	- block:
		- apt: name=apache2 state=present
		- template: src=httpd.conf.j2 dest=/etc/apache2/apache2.conf
		- service: name=apache2 state=started enabled=yes
	  when: ansible_os_family == 'Debian'
	  become: yes
```

Tasks inside the `block` will be run first. If there is a failure in any task in `block`, tasks inside `rescue` will be run. The tasks inside `always` will always be run, whether or not there were failures in either `block` or `rescue`.
It may not be necessary to use `block`/`rescue`/`always`.


