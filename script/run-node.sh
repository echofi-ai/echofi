
# remove existing daemon
rm -rf ~/.echofi/

# if $KEY exists it should be deleted
mnm_val="soap step crash ceiling path virtual this armor accident pond share track spice woman vault discover share holiday inquiry oak shine scrub bulb arrive"


# Set moniker and chain-id for Evmos (Moniker can be anything, chain-id must be an integer)
echofid init echofi-1 --chain-id testing-2 

echo "$mnm_val"| echofid keys add val --recover --keyring-backend test 

# Allocate genesis accounts (cosmos formatted addresses)
echofid genesis add-genesis-account val 100000000000000000000000000stake,1000000000000000000uecho --keyring-backend test

# Sign genesis transaction
echofid genesis gentx val 1000000000000000000000stake --keyring-backend test --chain-id testing-2

# Collect genesis tx
echofid genesis collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
echofid genesis validate-genesis

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
echofid start --pruning=nothing  --minimum-gas-prices=0.0001stake
