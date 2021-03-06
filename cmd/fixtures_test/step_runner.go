package fixturetest

import (
	"encoding/json"
	"strings"

	testing "github.com/Pylons-tech/pylons_sdk/cmd/fixtures_test/evtesting"

	inttest "github.com/Pylons-tech/pylons_sdk/cmd/test"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/msgs"

	"github.com/Pylons-tech/pylons_sdk/x/pylons/handlers"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RunMultiMsgTx is a function to send multiple messages in a transaction
// This support only 1 sender multi transaction for now
// TODO we need to support multi-message multi sender transaction
func RunMultiMsgTx(step FixtureStep, t *testing.T) {
	if len(step.MsgRefs) != 0 {
		var msgs []sdk.Msg
		var sender sdk.AccAddress
		for _, ref := range step.MsgRefs {
			var newMsg sdk.Msg
			switch ref.Action {
			case "fiat_item":
				msg := FiatItemMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
			case "update_item_string":
				msg := UpdateItemStringMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
			case "create_cookbook":
				msg := CreateCookbookMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
				t.Log("create_cookbook msg test", msg)
			case "create_recipe":
				msg := CreateRecipeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
			case "execute_recipe":
				msg := ExecuteRecipeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
			case "check_execution":
				msg := CheckExecutionMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
			case "create_trade":
				msg := CreateTradeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
			case "fulfill_trade":
				msg := FulfillTradeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
			case "disable_trade":
				msg := DisableTradeMsgFromRef(ref.ParamsRef, t)
				newMsg, sender = msg, msg.Sender
			}
			msgs = append(msgs, newMsg)
		}
		t.Log("sender", sender.String(), len(msgs), msgs, step.MsgRefs)
		txhash, err := inttest.SendMultiMsgTxWithNonce(t, msgs, sender.String(), true)
		t.MustNil(err)

		err = inttest.WaitForNextBlock()
		inttest.ErrValidation(t, "error waiting for check execution %+v", err)

		_, err = inttest.WaitAndGetTxData(txhash, inttest.GetMaxWaitBlock(), t)
		inttest.ErrValidation(t, "error getting tx result bytes %+v", err)

		CheckErrorOnTxFromTxHash(txhash, t)
		t.Log("txhash=", txhash)
	}
}

// CheckExecutionMsgFromRef collect check execution message from reference string
func CheckExecutionMsgFromRef(ref string, t *testing.T) msgs.MsgCheckExecution {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate execRef to execID
	newByteValue = UpdateExecID(newByteValue, t)

	var execType struct {
		ExecID        string
		PayToComplete bool
		Sender        sdk.AccAddress
	}
	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &execType)
	if err != nil {
		t.Fatal("error reading using GetAminoCdc execType", execType, err)
	}
	t.MustNil(err)

	return msgs.NewMsgCheckExecution(
		execType.ExecID,
		execType.PayToComplete,
		execType.Sender,
	)
}

// RunCheckExecution is a function to execute check execution
func RunCheckExecution(step FixtureStep, t *testing.T) {

	if step.ParamsRef != "" {
		chkExecMsg := CheckExecutionMsgFromRef(step.ParamsRef, t)
		txhash := inttest.TestTxWithMsgWithNonce(t, chkExecMsg, chkExecMsg.Sender.String(), true)

		err := inttest.WaitForNextBlock()
		inttest.ErrValidation(t, "error waiting for check execution %+v", err)

		txHandleResBytes, err := inttest.WaitAndGetTxData(txhash, inttest.GetMaxWaitBlock(), t)
		inttest.ErrValidation(t, "error getting tx result bytes %+v", err)

		CheckErrorOnTxFromTxHash(txhash, t)
		resp := handlers.CheckExecutionResp{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
		t.Log("txhash=", txhash)
		inttest.ErrValidation(t, "error unmarshaling tx response %+v", err)
		t.MustTrue(resp.Status == step.Output.TxResult.Status)
		if len(step.Output.TxResult.Message) > 0 {
			t.MustTrue(resp.Message == step.Output.TxResult.Message)
		}
	}
}

// FiatItemMsgFromRef collect check execution message from reference string
func FiatItemMsgFromRef(ref string, t *testing.T) msgs.MsgFiatItem {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate cookbook name to cookbook ID
	newByteValue = UpdateCBNameToID(newByteValue, t)

	var itemType types.Item
	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &itemType)
	if err != nil {
		t.Fatal("error reading using GetAminoCdc itemType", itemType, err)
	}
	t.MustNil(err)

	return msgs.NewMsgFiatItem(
		itemType.CookbookID,
		itemType.Doubles,
		itemType.Longs,
		itemType.Strings,
		itemType.Sender,
	)
}

// RunFiatItem is a function to execute fiat item
func RunFiatItem(step FixtureStep, t *testing.T) {

	if step.ParamsRef != "" {
		itmMsg := FiatItemMsgFromRef(step.ParamsRef, t)
		txhash := inttest.TestTxWithMsgWithNonce(t, itmMsg, itmMsg.Sender.String(), true)

		err := inttest.WaitForNextBlock()
		inttest.ErrValidation(t, "error waiting for fiat item %+v", err)

		txHandleResBytes, err := inttest.WaitAndGetTxData(txhash, inttest.GetMaxWaitBlock(), t)
		inttest.ErrValidation(t, "error getting tx result bytes %+v", err)

		CheckErrorOnTxFromTxHash(txhash, t)
		resp := handlers.FiatItemResponse{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)

		t.Log("txhash=", txhash)
		inttest.ErrValidation(t, "error unmarshaling tx response %+v", err)
		t.MustTrue(resp.ItemID != "")
	}
}

// UpdateItemStringMsgFromRef is a function to collect UpdateItemStringMsg from reference string
func UpdateItemStringMsgFromRef(ref string, t *testing.T) msgs.MsgUpdateItemString {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate item name to item ID
	newByteValue = UpdateItemIDFromName(newByteValue, false, t)

	var sTypeMsg msgs.MsgUpdateItemString
	err := json.Unmarshal(newByteValue, &sTypeMsg)
	if err != nil {
		t.Fatal("error reading using GetAminoCdc sTypeMsg", sTypeMsg, string(newByteValue), err)
	}
	t.MustNil(err)
	return sTypeMsg
}

// RunUpdateItemString is a function to update item's string value
func RunUpdateItemString(step FixtureStep, t *testing.T) {

	if step.ParamsRef != "" {
		sTypeMsg := UpdateItemStringMsgFromRef(step.ParamsRef, t)
		txhash := inttest.TestTxWithMsgWithNonce(t, sTypeMsg, sTypeMsg.Sender.String(), true)

		err := inttest.WaitForNextBlock()
		inttest.ErrValidation(t, "error waiting for set item field string %+v", err)

		txHandleResBytes, err := inttest.WaitAndGetTxData(txhash, inttest.GetMaxWaitBlock(), t)
		inttest.ErrValidation(t, "error getting tx result bytes %+v", err)

		CheckErrorOnTxFromTxHash(txhash, t)
		resp := handlers.UpdateItemStringResp{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)

		t.Log("txhash=", txhash)
		inttest.ErrValidation(t, "error unmarshaling tx response %+v", err)
	}
}

// CreateCookbookMsgFromRef is a function to get create cookbook message from reference
func CreateCookbookMsgFromRef(ref string, t *testing.T) msgs.MsgCreateCookbook {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)

	var cbType types.Cookbook
	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &cbType)
	if err != nil {
		t.Fatal("error reading using GetAminoCdc cbType", cbType, string(newByteValue), err)
	}
	t.MustNil(err)

	return msgs.NewMsgCreateCookbook(
		cbType.Name,
		cbType.ID,
		cbType.Description,
		cbType.Developer,
		cbType.Version,
		cbType.SupportEmail,
		cbType.Level,
		cbType.CostPerBlock,
		cbType.Sender,
	)
}

// RunCreateCookbook is a function to create cookbook
func RunCreateCookbook(step FixtureStep, t *testing.T) {
	if !FixtureTestOpts.CreateNewCookbook {
		return
	}
	if step.ParamsRef != "" {
		cbMsg := CreateCookbookMsgFromRef(step.ParamsRef, t)

		txhash := inttest.TestTxWithMsgWithNonce(t, cbMsg, cbMsg.Sender.String(), true)

		err := inttest.WaitForNextBlock()
		inttest.ErrValidation(t, "error waiting for creating cookbook %+v", err)

		txHandleResBytes, err := inttest.WaitAndGetTxData(txhash, inttest.GetMaxWaitBlock(), t)
		inttest.ErrValidationWithOutputLog(t, "error getting transaction data for creating cookbook %+v ----- %+v", txHandleResBytes, err)

		CheckErrorOnTxFromTxHash(txhash, t)
		resp := handlers.CreateCBResponse{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
		t.Log("txhash=", txhash)
		inttest.ErrValidation(t, "error unmarshaling tx response %+v", err)
		t.MustTrue(resp.CookbookID != "")
	}
}

// CreateRecipeMsgFromRef is a function to get create cookbook message from reference
func CreateRecipeMsgFromRef(ref string, t *testing.T) msgs.MsgCreateRecipe {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate cookbook name to cookbook id
	newByteValue = UpdateCBNameToID(newByteValue, t)
	// get item inputs from fileNames
	itemInputs := GetItemInputsFromBytes(newByteValue, t)
	// get entries from fileNames
	entries := GetEntriesFromBytes(newByteValue, t)

	var rcpTempl types.Recipe
	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &rcpTempl)
	if err != nil {
		t.Fatal("error reading using GetAminoCdc rcpTempl", rcpTempl, err)
	}
	t.MustNil(err)

	return msgs.NewMsgCreateRecipe(
		rcpTempl.Name,
		rcpTempl.CookbookID,
		rcpTempl.ID,
		rcpTempl.Description,
		rcpTempl.CoinInputs,
		itemInputs,
		entries,
		rcpTempl.Outputs,
		rcpTempl.BlockInterval,
		rcpTempl.Sender,
	)
}

// RunCreateRecipe is a function to create recipe
func RunCreateRecipe(step FixtureStep, t *testing.T) {
	if step.ParamsRef != "" {
		rcpMsg := CreateRecipeMsgFromRef(step.ParamsRef, t)

		txhash := inttest.TestTxWithMsgWithNonce(t, rcpMsg, rcpMsg.Sender.String(), true)

		err := inttest.WaitForNextBlock()
		inttest.ErrValidation(t, "error waiting for creating recipe %+v", err)

		txHandleResBytes, err := inttest.WaitAndGetTxData(txhash, inttest.GetMaxWaitBlock(), t)
		t.MustNil(err)

		CheckErrorOnTxFromTxHash(txhash, t)
		resp := handlers.CreateRecipeResponse{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
		t.Log("txhash=", txhash)
		inttest.ErrValidation(t, "error unmarshaling tx response %+v", err)
		t.MustTrue(resp.RecipeID != "")
	}
}

// ExecuteRecipeMsgFromRef collect execute recipe msg from reference string
func ExecuteRecipeMsgFromRef(ref string, t *testing.T) msgs.MsgExecuteRecipe {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate recipe name to recipe id
	newByteValue = UpdateRecipeName(newByteValue, t)
	// translate itemNames to itemIDs
	ItemIDs := GetItemIDsFromNames(newByteValue, false, t)

	var execType struct {
		RecipeID string
		Sender   sdk.AccAddress
		ItemIDs  []string `json:"ItemIDs"`
	}

	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &execType)
	if err != nil {
		t.Fatal("error reading using GetAminoCdc execType", execType, err)
	}
	t.MustNil(err)

	// t.Log("Executed recipe with below params", execType.RecipeID, execType.Sender, ItemIDs)
	return msgs.NewMsgExecuteRecipe(execType.RecipeID, execType.Sender, ItemIDs)
}

// RunExecuteRecipe is executed when an action "execute_recipe" is called
func RunExecuteRecipe(step FixtureStep, t *testing.T) {
	// TODO should check item ID is returned
	// TODO when items are generated, rather than returning whole should return only ID [if multiple, array of item IDs]

	if step.ParamsRef != "" {
		execMsg := ExecuteRecipeMsgFromRef(step.ParamsRef, t)
		txhash := inttest.TestTxWithMsgWithNonce(t, execMsg, execMsg.Sender.String(), true)

		err := inttest.WaitForNextBlock()
		inttest.ErrValidation(t, "error waiting for executing recipe %+v", err)

		if len(step.Output.TxResult.ErrorLog) > 0 {
			hmrErrMsg := inttest.GetHumanReadableErrorFromTxHash(txhash, t)
			t.Log("hmrErrMsg=", hmrErrMsg)
			t.MustTrue(strings.Contains(hmrErrMsg, step.Output.TxResult.ErrorLog))
		} else {
			txHandleResBytes, err := inttest.WaitAndGetTxData(txhash, inttest.GetMaxWaitBlock(), t)
			t.MustNil(err)
			CheckErrorOnTxFromTxHash(txhash, t)
			resp := handlers.ExecuteRecipeResp{}
			err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
			if err != nil {
				t.Fatal("failed to parse transaction result txhash=", txhash)
			}
			t.MustTrue(resp.Status == step.Output.TxResult.Status)
			if len(step.Output.TxResult.Message) > 0 {
				t.MustTrue(resp.Message == step.Output.TxResult.Message)
			}

			if resp.Message == "scheduled the recipe" { // delayed execution
				var scheduleRes handlers.ExecuteRecipeScheduleOutput

				err := json.Unmarshal(resp.Output, &scheduleRes)
				t.MustNil(err)
				execIDs[step.ID] = scheduleRes.ExecID
				for _, itemID := range execMsg.ItemIDs {
					item, err := inttest.GetItemByGUID(itemID)
					t.MustNil(err)
					t.MustTrue(len(item.OwnerRecipeID) != 0)
				}
				t.Log("scheduled execution", scheduleRes.ExecID)
			} else { // straight execution
				t.Log("straight execution result output", string(resp.Output))
			}
		}
	}
}

// CreateTradeMsgFromRef collect create trade msg from reference
func CreateTradeMsgFromRef(ref string, t *testing.T) msgs.MsgCreateTrade {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// get item inputs from fileNames
	itemInputs := GetItemInputsFromBytes(newByteValue, t)
	var trdType types.Trade
	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &trdType)
	if err != nil {
		t.Fatal("error reading using GetAminoCdc trdType", trdType, string(newByteValue), err)
	}
	t.MustTrue(err == nil)

	// get ItemOutputs from ItemOutputNames
	itemOutputs := GetItemOutputsFromBytes(newByteValue, trdType.Sender.String(), t)

	return msgs.NewMsgCreateTrade(
		trdType.CoinInputs,
		itemInputs,
		trdType.CoinOutputs,
		itemOutputs,
		trdType.ExtraInfo,
		trdType.Sender,
	)
}

// RunCreateTrade is a function to create trade
func RunCreateTrade(step FixtureStep, t *testing.T) {
	if step.ParamsRef != "" {
		createTrd := CreateTradeMsgFromRef(step.ParamsRef, t)
		t.Log("createTrd Msg=", createTrd)
		txhash := inttest.TestTxWithMsgWithNonce(t, createTrd, createTrd.Sender.String(), true)
		err := inttest.WaitForNextBlock()
		inttest.ErrValidation(t, "error while creating trade %+v", err)

		txHandleResBytes, err := inttest.WaitAndGetTxData(txhash, inttest.GetMaxWaitBlock(), t)
		inttest.ErrValidation(t, "error while waiting for create trade transaction %+v", err)

		CheckErrorOnTxFromTxHash(txhash, t)
		resp := handlers.CreateTradeResponse{}
		err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
		t.Log("txhash=", txhash)
		inttest.ErrValidation(t, "error unmarshaling tx response %+v", err)
		t.MustTrue(resp.TradeID != "")
	}
}

// FulfillTradeMsgFromRef collect fulfill trade message from reference string
func FulfillTradeMsgFromRef(ref string, t *testing.T) msgs.MsgFulfillTrade {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate extra info to trade id
	newByteValue = UpdateTradeExtraInfoToID(newByteValue, t)
	// translate itemNames to itemIDs
	ItemIDs := GetItemIDsFromNames(newByteValue, false, t)

	var trdType struct {
		TradeID string
		Sender  sdk.AccAddress
		ItemIDs []string `json:"ItemIDs"`
	}

	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &trdType)
	if err != nil {
		t.Fatal("error reading using GetAminoCdc trdType", trdType, err)
	}
	t.MustNil(err)

	return msgs.NewMsgFulfillTrade(trdType.TradeID, trdType.Sender, ItemIDs)
}

// RunFulfillTrade is a function to fulfill trade
func RunFulfillTrade(step FixtureStep, t *testing.T) {

	if step.ParamsRef != "" {
		ffTrdMsg := FulfillTradeMsgFromRef(step.ParamsRef, t)
		txhash := inttest.TestTxWithMsgWithNonce(t, ffTrdMsg, ffTrdMsg.Sender.String(), true)

		err := inttest.WaitForNextBlock()
		inttest.ErrValidation(t, "error waiting for fulfilling trade %+v", err)

		if len(step.Output.TxResult.ErrorLog) > 0 {
		} else {
			txHandleResBytes, err := inttest.WaitAndGetTxData(txhash, inttest.GetMaxWaitBlock(), t)
			t.MustNil(err)
			CheckErrorOnTxFromTxHash(txhash, t)
			resp := handlers.FulfillTradeResp{}
			err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
			if err != nil {
				t.Fatal("failed to parse transaction result txhash=", txhash)
			}
			t.MustTrue(resp.Status == step.Output.TxResult.Status)
			if len(step.Output.TxResult.Message) > 0 {
				t.MustTrue(resp.Message == step.Output.TxResult.Message)
			}
		}
	}
}

// DisableTradeMsgFromRef collect disable trade msg from reference string
func DisableTradeMsgFromRef(ref string, t *testing.T) msgs.MsgDisableTrade {
	byteValue := ReadFile(ref, t)
	// translate sender from account name to account address
	newByteValue := UpdateSenderKeyToAddress(byteValue, t)
	// translate extra info to trade id
	newByteValue = UpdateTradeExtraInfoToID(newByteValue, t)

	var trdType struct {
		TradeID string
		Sender  sdk.AccAddress
	}

	err := inttest.GetAminoCdc().UnmarshalJSON(newByteValue, &trdType)
	if err != nil {
		t.Fatal("error reading using GetAminoCdc trdType", trdType, err)
	}
	t.MustNil(err)

	return msgs.NewMsgDisableTrade(trdType.TradeID, trdType.Sender)
}

// RunDisableTrade is a function to disable trade
func RunDisableTrade(step FixtureStep, t *testing.T) {

	if step.ParamsRef != "" {
		dsTrdMsg := DisableTradeMsgFromRef(step.ParamsRef, t)
		txhash := inttest.TestTxWithMsgWithNonce(t, dsTrdMsg, dsTrdMsg.Sender.String(), true)

		err := inttest.WaitForNextBlock()
		inttest.ErrValidation(t, "error waiting for disabling trade %+v", err)

		if len(step.Output.TxResult.ErrorLog) > 0 {
		} else {
			txHandleResBytes, err := inttest.WaitAndGetTxData(txhash, inttest.GetMaxWaitBlock(), t)
			t.MustNil(err)
			CheckErrorOnTxFromTxHash(txhash, t)
			resp := handlers.DisableTradeResp{}
			err = inttest.GetAminoCdc().UnmarshalJSON(txHandleResBytes, &resp)
			if err != nil {
				t.Fatal("failed to parse transaction result txhash=", txhash)
			}
			t.MustTrue(resp.Status == step.Output.TxResult.Status)
			if len(step.Output.TxResult.Message) > 0 {
				t.MustTrue(resp.Message == step.Output.TxResult.Message)
			}
		}
	}
}
