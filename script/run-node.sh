
# remove existing daemon
rm -rf ~/.echofi/

# if $KEY exists it should be deleted
mnm_val="clip hire initial neck maid actor venue client foam budget lock catalog sweet steak waste crater broccoli pipe steak sister coyote moment obvious choose"

test="east night alone topic mom school foam stereo undo picnic runway print"
# Set moniker and chain-id for Evmos (Moniker can be anything, chain-id must be an integer)
echofid init echofi-1 --chain-id testing-2 

echo "$mnm_val"| echofid keys add val --recover --keyring-backend test 
echo "$test"| echofid keys add test --recover --keyring-backend test 

# Allocate genesis accounts (cosmos formatted addresses)
echofid genesis add-genesis-account val 100000000000000000000000000stake,10000000000000000000000uecho --keyring-backend test

# Sign genesis transaction
echofid genesis gentx val 1000000000000000000000uecho --keyring-backend test --chain-id testing-2

# Collect genesis tx
echofid genesis collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
echofid genesis validate-genesis

update_test_genesis () {
    # EX: update_test_genesis '.consensus_params["block"]["max_gas"]="100000000"'
    cat $HOME/.echofi/config/genesis.json | jq "$1" > tmp.json && mv tmp.json $HOME/.echofi/config/genesis.json
}

# update_test_genesis '.app_state["gov"]["voting_params"]["voting_period"] = "30s"'
update_test_genesis '.app_state["mint"]["params"]["mint_denom"]= "uecho"'
update_test_genesis '.app_state["gov"]["params"]["min_deposit"]=[{"denom": "uecho","amount": "1000000"}]'
update_test_genesis '.app_state["crisis"]["constant_fee"]={"denom": "uecho","amount": "1000"}'
update_test_genesis '.app_state["staking"]["params"]["bond_denom"]="uecho"'

# # validator2
# VALIDATOR2_CONFIG=$HOME/.echofi/config/config.toml
# sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26655|g' $VALIDATOR2_CONFIG
# sed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26654|g' $VALIDATOR2_CONFIG
# sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26653|g' $VALIDATOR2_CONFIG
# sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $VALIDATOR2_CONFIG
# sed -i -E 's|prometheus = false|prometheus = true|g' $VALIDATOR2_CONFIG
# sed -i -E 's|prometheus_listen_addr = ":26660"|prometheus_listen_addr = ":26630"|g' $VALIDATOR2_CONFIG

# VALIDATOR2_APP_TOML=$HOME/.echofi/config/app.toml
# sed -i -E 's|tcp://localhost:1317|tcp://localhost:1316|g' $VALIDATOR2_APP_TOML
# sed -i -E 's|localhost:9090|localhost:9088|g' $VALIDATOR2_APP_TOML
# sed -i -E 's|localhost:9091|localhost:9089|g' $VALIDATOR2_APP_TOML
# sed -i -E 's|tcp://0.0.0.0:10337|tcp://0.0.0.0:10347|g' $VALIDATOR2_APP_TOML

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
echofid start --pruning=nothing  --minimum-gas-prices=0.0001uecho
