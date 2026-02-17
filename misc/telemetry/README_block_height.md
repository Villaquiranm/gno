# https://github.com/gnolang/gno/security/advisories/GHSA-4m6f-9qpf-4qjv

This Poc modifies the gnolang command in order to try to validate this security issue:

To do that you need to go to the root of the project and execute
```
docker build . -t gnolang:local -f Dockerfile.gnoland
```
This will build your own local gnolang image with the modifications (vulnerabilities) that we induced on the code.

the malicious node includes a new rpc method /mode
`http://localhost:26657/mode?mode=malicious`
if you go to this address you'll set the node as malicious it will send a fake height and then return to the normal as soon as the fake height is broadcasted


be sure to remove all volumes before rexecuting the docker compose 
```sh
docker compose down
docker volume prune -a
```




