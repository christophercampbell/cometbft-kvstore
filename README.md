# cometbft-kvstore

kvstore with cometbft

https://docs.cometbft.com/v0.37/guides/go-built-in

```
curl -s 'localhost:26657/broadcast_tx_commit?tx="cometbft=rocks"'
```

```
 curl -s 'localhost:26657/abci_query?data="cometbft"' | jq '.result.response.value' | xargs -I{} bash -c 'base64 -d <<< {}'
```


Setup 4 nodes as persistent peers

```
cometbft init --home <pathN>
```

get peer node ids

```
cometbft show-node-id --home <pathN>
```

- `<pathN>/config/genesis.toml`: configure genesis with each validator's public info and voting power, and *same chain ID*
- `<pathN>/config/config.toml`: add all nodes as persistent_peers: `node-id@ip:port` in CSV format

Start the servers 

```
kvstore --cmt-home <pathN>
```
