# What is it?
```
$ docker run --rm clovergrp/vault-copy -h
Utility to copy with replace whole branches in vault
Usage: vault-copy [Options]
Options:
  -i string
    	Path to copy
  -o string
    	Path where to copy
  -p int
    	Password length (default 15)
  -r string
    	Sed regular expression to replace old variables (see https://github.com/rwtodd/Go.Sed)
  -t string
    	Path to file with token (default "./token")
```
# How to install
```
wget https://github.com/Clover-Group/vault-copy/releases/download/1.0.0/vault-copy_1.0.0_linux_amd64
chmod +x vault-copy_1.0.0_linux_amd64
sudo mv vault-copy_1.0.0_linux_amd64 /usr/local/bin/vault-copy
```
Or you can just pull image from docker hub:
```
docker pull clovergrp/vault-copy
```
# How to build
```
$ docker build -t clovergrp/vault-copy .
```
# How to run
```
$ docker run -v ~/token:/token -e VAULT_ADDR=http://localhost:8200 --rm clovergrp/vault-copy -i old_branch -o new_branch -t /token
```
