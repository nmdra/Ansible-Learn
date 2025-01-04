1.`-f 1` to tell Ansible to use only one fork.
```bash
$ ansible multi -a "hostname" -f 1
192.168.56.4 | CHANGED | rc=0 >>
orc-app1.test
192.168.56.5 | CHANGED | rc=0 >>
orc-app2.test
192.168.56.6 | CHANGED | rc=0 >>
orc-db.test
```

2. Install chrony service (https://chrony-project.org/)
```bash
ansible multi -b -m dnf -a "name=chrony state=present"
```

> [!TIP]
> **Chrony** is an open-source software suite designed for maintaining accurate system time in Linux-based and UNIX-like operating systems.


