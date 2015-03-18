# Quickstart guide

The quickest way to install hello is to download the latest release:

```shell
# download the latest release
export HELLO_VERSION=v0.0.1-alpha.1
cd /tmp
wget https://github.com/gravitational/hello/releases/download/$VERSION/hello-$HELLO_VERSION-linux-amd64.tar.gz
tar -xzf hello-$HELLO_VERSION-linux-amd64.tar.gz

# start a hello server listening on localhost:23456 
# and using etcd as a configuration backend
/tmp/hello -addr=localhost:23456\
           -backend=etcd\
           -backendConfig='{
              "nodes": ["http://localhost:4001"], 
              "key": "/hello"}'

# upsert a greeting
/tmp/hctl -hello=http://localhost:23456 greeting upsert -id=hello.us -val=Hello

# execute a hello request
/tmp/hctl -hello=http://localhost:23456 hello -id=hello.us -name=Dog
```
