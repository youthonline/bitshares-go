package database

import (
	"encoding/json"
	"github.com/scorum/openledger-go/caller"
	"github.com/scorum/openledger-go/types"
)

type API struct {
	caller caller.Caller
	id     caller.APIID
}

func NewAPI(id caller.APIID, caller caller.Caller) *API {
	return &API{id: id, caller: caller}
}

func (api *API) call(method string, args []interface{}, reply interface{}) error {
	return api.caller.Call(api.id, method, args, reply)
}

func (api *API) setCallback(method string, callback func(raw json.RawMessage)) error {
	return api.caller.SetCallback(api.id, method, callback)
}

func (api *API) GetChainID() (*string, error) {
	var resp string
	err := api.call("get_chain_id", caller.EmptyParams, &resp)
	return &resp, err
}

// GetConfig retrieves compile-time constants
func (api *API) GetConfig() (*Config, error) {
	var config Config
	err := api.call("get_config", caller.EmptyParams, &config)
	return &config, err
}

// GetDynamicGlobalProperties retrieves the current global_property_object
func (api *API) GetDynamicGlobalProperties() (*DynamicGlobalProperties, error) {
	var resp DynamicGlobalProperties
	err := api.call("get_dynamic_global_properties", caller.EmptyParams, &resp)
	return &resp, err
}

// LookupAssetSymbols get assets corresponding to the provided symbols or IDs
func (api *API) LookupAssetSymbols(symbols ...string) ([]*Asset, error) {
	var resp []*Asset
	err := api.call("lookup_asset_symbols", []interface{}{symbols}, &resp)
	return resp, err
}

// GetLimitOrders returns limit orders in a given market.
// There are both sell and buy orders.
// For the sell orders LimitOrder.SellPrice.Base = the given base
// For the buy orders LimitOrder.SellPrice.Base = the given quote
func (api *API) GetLimitOrders(base, quote types.ObjectID, limit uint32) ([]*LimitOrder, error) {
	var resp []*LimitOrder
	err := api.call("get_limit_orders", []interface{}{base.String(), quote.String(), limit}, &resp)
	return resp, err
}

// GetBlockHeader returns block header by the given block number
func (api *API) GetBlockHeader(blockNum int32) (*BlockHeader, error) {
	var resp BlockHeader
	err := api.call("get_block_header", []interface{}{blockNum}, &resp)
	return &resp, err
}

// GetBlock
func (api *API) GetBlock(blockNum int32) (*Block, error) {
	var resp Block
	err := api.call("get_block", []interface{}{blockNum}, &resp)
	return &resp, err
}

// GetTicker returns the ticker for the market assetA:assetB (past 24 hours)
func (api *API) GetTicker(base, quote types.ObjectID) (*MarketTicker, error) {
	var resp MarketTicker
	err := api.call("get_ticker", []interface{}{base.String(), quote.String()}, &resp)
	return &resp, err
}

// SetBlockAppliedCallback registers a global subscription callback
func (api *API) SetBlockAppliedCallback(notice func(blockID string, err error)) (err error) {
	err = api.setCallback("set_block_applied_callback", func(raw json.RawMessage) {
		var header []string
		if err := json.Unmarshal(raw, &header); err != nil {
			notice("", err)
		}
		for _, b := range header {
			notice(b, nil)
		}
	})
	return
}

// CancelAllSubscriptions cancel all subscriptions
func (api *API) CancelAllSubscriptions() error {
	return api.call("cancel_all_subscriptions", caller.EmptyParams, nil)
}
