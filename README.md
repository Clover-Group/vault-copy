# How to build
```
docker build -t clovergrp/vault-copy .
```
# How to run
```
docker run -v ~/token:/token -e VAULT_ADDR=http://localhost:8200 --rm clovergrp/vault-copy -i old_branch -o new_branch -t /token
```
