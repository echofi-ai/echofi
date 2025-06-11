#!/bin/bash
set -xeu

# always returns true so set -e doesn't exit if it is not running.
killall echofid || true
rm -rf $HOME/.echofid/

# make four echofi directories
mkdir $HOME/.echofid
cd $HOME/.echofid/
mkdir $HOME/.echofid/validator1
mkdir $HOME/.echofid/validator2
mkdir $HOME/.echofid/validator3

# init all three validators
echofid init --chain-id=testing-1 validator1 --home=$HOME/.echofid/validator1
echofid init --chain-id=testing-1 validator2 --home=$HOME/.echofid/validator2
echofid init --chain-id=testing-1 validator3 --home=$HOME/.echofid/validator3

# create keys for all three validators
# cosmos1f7twgcq4ypzg7y24wuywy06xmdet8pc4473tnq
echo $(cat /Users/donglieu/script/keys/mnemonic1)| echofid keys add validator1 --recover --keyring-backend=test --home=$HOME/.echofid/validator1
# cosmos1w7f3xx7e75p4l7qdym5msqem9rd4dyc4752spg
echo $(cat /Users/donglieu/script/keys/mnemonic2)| echofid keys add validator2 --recover --keyring-backend=test --home=$HOME/.echofid/validator2
# cosmos1g9v3zjt6rfkwm4s8sw9wu4jgz9me8pn27f8nyc
echo $(cat /Users/donglieu/script/keys/mnemonic3)| echofid keys add validator3 --recover --keyring-backend=test --home=$HOME/.echofid/validator3

# create validator node with tokens to transfer to the three other nodes
echofid genesis add-genesis-account $(echofid keys show validator1 -a --keyring-backend=test --home=$HOME/.echofid/validator1) 10000000000000000000000000000000stake,10000000000000000000000000000000uecho --home=$HOME/.echofid/validator1 
echofid genesis add-genesis-account $(echofid keys show validator2 -a --keyring-backend=test --home=$HOME/.echofid/validator2) 10000000000000000000000000000000stake,10000000000000000000000000000000uecho --home=$HOME/.echofid/validator1 
echofid genesis add-genesis-account $(echofid keys show validator3 -a --keyring-backend=test --home=$HOME/.echofid/validator3) 10000000000000000000000000000000stake,10000000000000000000000000000000uecho --home=$HOME/.echofid/validator1
echofid genesis add-genesis-account $(echofid keys show validator1 -a --keyring-backend=test --home=$HOME/.echofid/validator1) 10000000000000000000000000000000stake,10000000000000000000000000000000uecho --home=$HOME/.echofid/validator2
echofid genesis add-genesis-account $(echofid keys show validator2 -a --keyring-backend=test --home=$HOME/.echofid/validator2) 10000000000000000000000000000000stake,10000000000000000000000000000000uecho --home=$HOME/.echofid/validator2 
echofid genesis add-genesis-account $(echofid keys show validator3 -a --keyring-backend=test --home=$HOME/.echofid/validator3) 10000000000000000000000000000000stake,10000000000000000000000000000000uecho --home=$HOME/.echofid/validator2 
echofid genesis add-genesis-account $(echofid keys show validator1 -a --keyring-backend=test --home=$HOME/.echofid/validator1) 10000000000000000000000000000000stake,10000000000000000000000000000000uecho --home=$HOME/.echofid/validator3 
echofid genesis add-genesis-account $(echofid keys show validator2 -a --keyring-backend=test --home=$HOME/.echofid/validator2) 10000000000000000000000000000000stake,10000000000000000000000000000000uecho --home=$HOME/.echofid/validator3 
echofid genesis add-genesis-account $(echofid keys show validator3 -a --keyring-backend=test --home=$HOME/.echofid/validator3) 10000000000000000000000000000000stake,10000000000000000000000000000000uecho --home=$HOME/.echofid/validator3
echofid genesis gentx validator1 1000000000000000000000stake --keyring-backend=test --home=$HOME/.echofid/validator1 --chain-id=testing-1
echofid genesis gentx validator2 1000000000000000000000stake --keyring-backend=test --home=$HOME/.echofid/validator2 --chain-id=testing-1
echofid genesis gentx validator3 1000000000000000000000stake --keyring-backend=test --home=$HOME/.echofid/validator3 --chain-id=testing-1

cp validator2/config/gentx/*.json $HOME/.echofid/validator1/config/gentx/
cp validator3/config/gentx/*.json $HOME/.echofid/validator1/config/gentx/
echofid genesis collect-gentxs --home=$HOME/.echofid/validator1 
echofid genesis collect-gentxs --home=$HOME/.echofid/validator2
echofid genesis collect-gentxs --home=$HOME/.echofid/validator3 

cp validator1/config/genesis.json $HOME/.echofid/validator2/config/genesis.json
cp validator1/config/genesis.json $HOME/.echofid/validator3/config/genesis.json


# change app.toml values
VALIDATOR1_APP_TOML=$HOME/.echofid/validator1/config/app.toml
VALIDATOR2_APP_TOML=$HOME/.echofid/validator2/config/app.toml
VALIDATOR3_APP_TOML=$HOME/.echofid/validator3/config/app.toml

# validator1
sed -i -E 's|localhost:9090|localhost:9050|g' $VALIDATOR1_APP_TOML
sed -i -E 's|minimum-gas-prices = ""|minimum-gas-prices = "0.0001stake"|g' $VALIDATOR1_APP_TOML

# validator2
sed -i -E 's|tcp://localhost:1317|tcp://localhost:1316|g' $VALIDATOR2_APP_TOML
# sed -i -E 's|0.0.0.0:9090|0.0.0.0:9088|g' $VALIDATOR2_APP_TOML
sed -i -E 's|localhost:9090|localhost:9088|g' $VALIDATOR2_APP_TOML
# sed -i -E 's|0.0.0.0:9091|0.0.0.0:9089|g' $VALIDATOR2_APP_TOML
sed -i -E 's|localhost:9091|localhost:9089|g' $VALIDATOR2_APP_TOML
sed -i -E 's|minimum-gas-prices = ""|minimum-gas-prices = "0.0001stake"|g' $VALIDATOR2_APP_TOML

# validator3
sed -i -E 's|tcp://localhost:1317|tcp://localhost:1315|g' $VALIDATOR3_APP_TOML
# sed -i -E 's|0.0.0.0:9090|0.0.0.0:9086|g' $VALIDATOR3_APP_TOML
sed -i -E 's|localhost:9090|localhost:9086|g' $VALIDATOR3_APP_TOML
# sed -i -E 's|0.0.0.0:9091|0.0.0.0:9087|g' $VALIDATOR3_APP_TOML
sed -i -E 's|localhost:9091|localhost:9087|g' $VALIDATOR3_APP_TOML
sed -i -E 's|minimum-gas-prices = ""|minimum-gas-prices = "0.0001stake"|g' $VALIDATOR3_APP_TOML

# change config.toml values
VALIDATOR1_CONFIG=$HOME/.echofid/validator1/config/config.toml
VALIDATOR2_CONFIG=$HOME/.echofid/validator2/config/config.toml
VALIDATOR3_CONFIG=$HOME/.echofid/validator3/config/config.toml


# validator1
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR1_CONFIG
# sed -i -E 's|prometheus = false|prometheus = true|g' $VALIDATOR1_CONFIG


# validator2
sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26655|g' $VALIDATOR2_CONFIG
sed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26654|g' $VALIDATOR2_CONFIG
sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26653|g' $VALIDATOR2_CONFIG
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR2_CONFIG
sed -i -E 's|prometheus = false|prometheus = true|g' $VALIDATOR2_CONFIG
sed -i -E 's|prometheus_listen_addr = ":26660"|prometheus_listen_addr = ":26630"|g' $VALIDATOR2_CONFIG

# validator3
sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26652|g' $VALIDATOR3_CONFIG
sed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26651|g' $VALIDATOR3_CONFIG
sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26650|g' $VALIDATOR3_CONFIG
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR3_CONFIG
sed -i -E 's|prometheus = false|prometheus = true|g' $VALIDATOR3_CONFIG
sed -i -E 's|prometheus_listen_addr = ":26660"|prometheus_listen_addr = ":26620"|g' $VALIDATOR3_CONFIG

# # update
# update_test_genesis () {
#     # EX: update_test_genesis '.consensus_params["block"]["max_gas"]="100000000"'
#     cat $HOME/.echofid/validator1/config/genesis.json | jq "$1" > tmp.json && mv tmp.json $HOME/.echofid/validator1/config/genesis.json
#     cat $HOME/.echofid/validator2/config/genesis.json | jq "$1" > tmp.json && mv tmp.json $HOME/.echofid/validator2/config/genesis.json
#     cat $HOME/.echofid/validator3/config/genesis.json | jq "$1" > tmp.json && mv tmp.json $HOME/.echofid/validator3/config/genesis.json
# }

# update_test_genesis '.app_state["gov"]["voting_params"]["voting_period"] = "30s"'
# update_test_genesis '.app_state["mint"]["params"]["mint_denom"]= "stake"'
# update_test_genesis '.app_state["gov"]["deposit_params"]["min_deposit"]=[{"denom": "stake","amount": "1000000"}]'
# update_test_genesis '.app_state["crisis"]["constant_fee"]={"denom": "stake","amount": "1000"}'
# update_test_genesis '.app_state["staking"]["params"]["bond_denom"]="stake"'


# copy validator1 genesis file to validator2-3
cp $HOME/.echofid/validator1/config/genesis.json $HOME/.echofid/validator2/config/genesis.json
cp $HOME/.echofid/validator1/config/genesis.json $HOME/.echofid/validator3/config/genesis.json

# copy tendermint node id of validator1 to persistent peers of validator2-3
node1=$(echofid tendermint show-node-id --home=$HOME/.echofid/validator1)
node2=$(echofid tendermint show-node-id --home=$HOME/.echofid/validator2)
node3=$(echofid tendermint show-node-id --home=$HOME/.echofid/validator3)
sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$node1@localhost:26656,$node2@localhost:26653,$node3@localhost:26650\"|g" $HOME/.echofid/validator1/config/config.toml
sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$node1@localhost:26656,$node2@localhost:26653,$node3@localhost:26650\"|g" $HOME/.echofid/validator2/config/config.toml
sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"$node1@localhost:26656,$node2@localhost:26653,$node3@localhost:26650\"|g" $HOME/.echofid/validator3/config/config.toml


# # start all three validators/
# echofid start --home=$HOME/.echofid/validator1
screen -S echofi1 -t echofi1 -d -m echofid start --home=$HOME/.echofid/validator1
screen -S echofi2 -t echofi2 -d -m echofid start --home=$HOME/.echofid/validator2
screen -S echofi3 -t echofi3 -d -m echofid start --home=$HOME/.echofid/validator3
# echofid start --home=$HOME/.echofid/validator3

# screen -r echofi1

sleep 7

echofid tx bank send $(echofid keys show validator1 -a --keyring-backend=test --home=$HOME/.echofid/validator1) $(echofid keys show validator2 -a --keyring-backend=test --home=$HOME/.echofid/validator2) 100000stake --keyring-backend=test --chain-id=testing-1 -y --home=$HOME/.echofid/validator1 --fees 10stake

# sleep 7
# echofid tx bank send cosmos1f7twgcq4ypzg7y24wuywy06xmdet8pc4473tnq cosmos1qvuhm5m644660nd8377d6l7yz9e9hhm9evmx3x 10000000000000000000000stake --keyring-backend=test --chain-id=testing-1 -y --home=$HOME/.echofid/validator1 --fees 200000stake
# sleep 7
# echofid tx bank send echofi1f7twgcq4ypzg7y24wuywy06xmdet8pc4hhtf9t echofi16gjg8p5fedy48wf403jwmz2cxlwqtkqlwe0lug 10000000000000000000000stake --keyring-backend=test --chain-id=testing-1 -y --home=$HOME/.echofid/validator1 --fees 10stake

