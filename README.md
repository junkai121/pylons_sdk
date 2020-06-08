# pylons SDK

pylons SDK provides packages to build blockchain games on pylons eco system.

## Setup development environment

```
git clone https://github.com/Pylons-tech/pylons_sdk
brew install pre-commit
brew install golangci/tap/golangci-lint
go get -u golang.org/x/lint/golint
pre-commit install
```

# SDK publish preparation

Check below things before publishing
```
make fixture_tests
```

# How to add feature

All the features added should have fixture test and it should be well documented.


# Packages

## Fixture Test Package
"github.com/Pylons-tech/pylons_sdk/cmd/fixtures_test"

## Integration Test Util Package
"github.com/Pylons-tech/pylons_sdk/cmd/test"

- GetAccountAddr
- GetAccountInfoFromName
- ListItemsViaCLI
- GetDaemonStatus
- CLIOpts.CustomNode
- CLIOpts.RestEndpoint
- WaitAndGetTxData
- ReadFile
- GetAminoCdc
- RunPylonsCli
- CleanFile
- GenTxWithMsg
- WaitForNextBlock
- WaitAndGetTxData
- GetHumanReadableErrorFromTxHash
- TestTxWithMsgWithNonce
- GetItemByGUID

- SendMultiMsgTxWithNonce
SendMultiMsgTxWithNonce is an integration test utility to send multiple message transaction from a single sender, single signed transaction.



"github.com/Pylons-tech/pylons_sdk/x/pylons/handlers"

Structs

- handlers.ExecuteRecipeResp
- handlers.ExecuteRecipeScheduleOutput{}
- handlers.CheckExecutionResp{}
- handlers.CreateCBResponse{}
- handlers.CreateRecipeResponse{}
- handlers.FulfillTradeResp{}
- handlers.PopularRecipeType
- handlers.GetParamsForPopularRecipe
- handlers.FiatItemResponse{}
- handlers.UpdateItemStringResp{}

"github.com/Pylons-tech/pylons_sdk/x/pylons/msgs"

All msg types
- MsgCheckExecution
- MsgCreateCookbook
- MsgCreateRecipe
- MsgCreateTrade
- MsgDisableRecipe
- MsgDisableTrade
- MsgEnableRecipe
- MsgEnableTrade
- MsgExecuteRecipe
- MsgFiatItem
- MsgFulfillTrade
- MsgGetPylons
- MsgSendPylons
- MsgUpdateItemString
- MsgUpdateCookbook
- MsgUpdateRecipe

Utility functions 

- msgs.NewMsgGetPylons
- msgs.NewMsgExecuteRecipe
- msgs.NewMsgCreateCookbook
- msgs.NewMsgGetPylons
- msgs.NewMsgUpdateItemString
- msgs.NewMsgCreateTrade
- msgs.NewMsgFulfillTrade
- msgs.NewMsgDisableTrade
- msgs.NewMsgCheckExecution 
- msgs.NewMsgFiatItem
- msgs.NewMsgCreateRecipe
- msgs.DefaultCostPerBlock


"github.com/Pylons-tech/pylons_sdk/x/pylons/types"


structs

- types.Item
- types.Cookbook
- types.Recipe
- types.Trade
- types.FloatString
- types.EntriesList
- types.TradeList
- types.Execution
- types.CoinOutput
- types.ItemModifyParams
- types.PremiumTier.Fee
- types.ItemList 
- types.ItemInputList
- types.ItemInput
- types.DoubleInputParamList
- types.DoubleInputParam
- types.LongInputParamList
- types.LongInputParam
- types.StringInputParamList
- types.StringInputParam
- types.CoinInputList,
- types.WeightedOutputsList,

Utility functions 

- types.NewPylon
- types.GenItemInputList
- types.GenEntries
- types.GenCoinInputList
- types.ItemInputList{},
- types.GenItemOnlyEntry
- types.GenCoinInputList
- types.GenEntriesFirstItemNameUpgrade(desItemName),
- types.GenOneOutput

"github.com/Pylons-tech/pylons_sdk/app"

- app.MakeCodec()

"github.com/Pylons-tech/pylons_sdk/x/pylons/queriers"

- queriers.ExecResp
- queriers.ItemResp
