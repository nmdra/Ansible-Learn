# Ansible on Host Machine

In Ansible, **localhost** is automatically included in the inventory by default.
This means you can run tasks on the control machine (the machine where Ansible is executed) without needing to explicitly define it in the inventory file.

- By default, Ansible uses the **`local` connection plugin** for `localhost`, bypassing SSH. This makes it efficient for local operations.
- Example configuration (inferred automatically):
```ini
[localhost] 
localhost ansible_connection=local
```
