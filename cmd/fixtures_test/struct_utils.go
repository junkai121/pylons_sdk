package fixturetest

import (
	"encoding/json"
	"io/ioutil"
	"os"

	testing "github.com/Pylons-tech/pylons_sdk/cmd/fixtures_test/evtesting"

	inttest "github.com/Pylons-tech/pylons_sdk/cmd/test"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
)

var execIDs map[string]string = make(map[string]string)

// ReadFile is a function to read file
func ReadFile(fileURL string, t *testing.T) []byte {
	jsonFile, err := os.Open(fileURL)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}

// UnmarshalIntoEmptyInterface is a function to convert bytes into map[string]interface{}
func UnmarshalIntoEmptyInterface(bytes []byte, t *testing.T) map[string]interface{} {
	var raw map[string]interface{}
	if err := json.Unmarshal(bytes, &raw); err != nil {
		t.Fatal("read raw file using json.Unmarshal:", err)
	}
	return raw
}

// UpdateSenderKeyToAddress is a function to update sender key to sender's address
func UpdateSenderKeyToAddress(bytes []byte, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	senderName, ok := raw["Sender"].(string)
	t.MustTrue(ok)
	raw["Sender"] = inttest.GetAccountAddr(senderName, t)
	newBytes, err := json.Marshal(raw)
	t.MustNil(err)
	return newBytes
}

// UpdateCBNameToID is a function to update cookbook name to cookbook id if it has cookbook name field
func UpdateCBNameToID(bytes []byte, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	cbName, ok := raw["CookbookName"].(string)
	if !ok {
		return bytes
	}
	cbID, exist, err := inttest.GetCookbookIDFromName(cbName, "")
	if exist && err != nil {
		raw["CookbookID"] = cbID
		newBytes, err := json.Marshal(raw)
		t.MustNil(err)
		return newBytes
	}
	return bytes
}

// UpdateRecipeName is a function to update recipe name into recipe id
func UpdateRecipeName(bytes []byte, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	rcpName, ok := raw["RecipeName"].(string)
	t.MustTrue(ok)
	rcpID, exist, err := inttest.GetRecipeIDFromName(rcpName)
	t.MustTrue(exist)
	t.MustNil(err)
	raw["RecipeID"] = rcpID
	newBytes, err := json.Marshal(raw)
	t.MustNil(err)
	return newBytes
}

// UpdateTradeExtraInfoToID is a function to update trade extra info into trade id
func UpdateTradeExtraInfoToID(bytes []byte, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	trdInfo, ok := raw["TradeInfo"].(string)
	t.MustTrue(ok)
	trdID, exist, err := inttest.GetTradeIDFromExtraInfo(trdInfo)
	t.MustTrue(exist)
	t.MustNil(err)
	raw["TradeID"] = trdID
	newBytes, err := json.Marshal(raw)
	t.MustNil(err)
	return newBytes
}

// UpdateExecID is a function to set execute id from execID reference
func UpdateExecID(bytes []byte, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	var execRefReader struct {
		ExecRef string
	}
	if err := json.Unmarshal(bytes, &execRefReader); err != nil {
		t.Fatal("read execRef using json.Unmarshal:", err)
	}

	var ok bool
	raw["ExecID"], ok = execIDs[execRefReader.ExecRef]
	if !ok {
		t.Fatal("execID not available for ref=", execRefReader.ExecRef)
	}
	newBytes, err := json.Marshal(raw)
	t.MustNil(err)
	return newBytes
}

// UpdateItemIDFromName is a function to set item id from item name
func UpdateItemIDFromName(bytes []byte, includeLockedByRcp bool, t *testing.T) []byte {
	raw := UnmarshalIntoEmptyInterface(bytes, t)

	itemName, ok := raw["ItemName"].(string)

	t.MustTrue(ok)
	itemID, exist, err := inttest.GetItemIDFromName(itemName, includeLockedByRcp)
	if !exist {
		t.Log("no item named=", itemName, "and includeLockedByRcp=", includeLockedByRcp)
	}
	t.MustTrue(exist)
	t.MustNil(err)
	raw["ItemID"] = itemID
	newBytes, err := json.Marshal(raw)
	t.MustNil(err)
	return newBytes
}

// GetItemIDsFromNames is a function to set item ids from names for recipe execution
func GetItemIDsFromNames(bytes []byte, includeLockedByRcp bool, t *testing.T) []string {
	var itemNamesResp struct {
		ItemNames []string
	}
	if err := json.Unmarshal(bytes, &itemNamesResp); err != nil {
		t.Fatal("read item names using json.Unmarshal:", err)
	}
	ItemIDs := []string{}

	for _, itemName := range itemNamesResp.ItemNames {
		itemID, exist, err := inttest.GetItemIDFromName(itemName, includeLockedByRcp)
		if !exist {
			t.Log("no item named=", itemName, "and includeLockedByRcp=", includeLockedByRcp)
		}
		t.MustTrue(exist)
		t.MustNil(err)
		ItemIDs = append(ItemIDs, itemID)
	}
	return ItemIDs
}

// GetItemInputsFromBytes is a function to get item input list from bytes
func GetItemInputsFromBytes(bytes []byte, t *testing.T) types.ItemInputList {
	var itemInputRefsReader struct {
		ItemInputRefs []string
	}
	if err := json.Unmarshal(bytes, &itemInputRefsReader); err != nil {
		t.Fatal("read itemInputRefsReader using json.Unmarshal:", err)
	}

	var itemInputs types.ItemInputList

	for _, iiRef := range itemInputRefsReader.ItemInputRefs {
		var ii types.ItemInput
		iiBytes := ReadFile(iiRef, t)
		err := inttest.GetAminoCdc().UnmarshalJSON(iiBytes, &ii)
		if err != nil {
			t.Fatal("error parsing item input provided via fixture error=", err)
		}
		itemInputs = append(itemInputs, ii)
	}
	return itemInputs
}

// GetItemOutputsFromBytes is a function to get item outputs from bytes
func GetItemOutputsFromBytes(bytes []byte, sender string, t *testing.T) types.ItemList {
	var itemOutputNamesReader struct {
		ItemOutputNames []string
	}
	if err := json.Unmarshal(bytes, &itemOutputNamesReader); err != nil {
		t.Fatal("read itemOutputNamesReader using json.Unmarshal:", err)
	}

	var itemOutputs types.ItemList

	for _, iN := range itemOutputNamesReader.ItemOutputNames {
		var io types.Item
		iID, ok, err := inttest.GetItemIDFromName(iN, false)
		t.MustNil(err)
		t.MustTrue(ok)
		io, err = inttest.GetItemByGUID(iID)
		t.MustNil(err)
		itemOutputs = append(itemOutputs, io)
	}
	return itemOutputs
}

// GetEntriesFromBytes is a function to get entries from bytes
func GetEntriesFromBytes(bytes []byte, t *testing.T) types.EntriesList {
	var entriesReader struct {
		Entries struct {
			CoinOutputs []types.CoinOutput
			ItemOutputs []struct {
				Ref        string
				ModifyItem struct {
					ItemInputRef    *int
					ModifyParamsRef string
				}
			}
		}
	}

	if err := json.Unmarshal(bytes, &entriesReader); err != nil {
		t.Fatal("read entriesReader using json.Unmarshal:", err)
	}

	var wpl types.EntriesList
	for _, co := range entriesReader.Entries.CoinOutputs {
		wpl = append(wpl, co)
	}

	for _, io := range entriesReader.Entries.ItemOutputs {
		var pio types.ItemOutput
		if len(io.Ref) > 0 {
			ioBytes := ReadFile(io.Ref, t)
			err := json.Unmarshal(ioBytes, &pio)
			if err != nil {
				t.Fatal("error parsing item output provided via fixture Bytes=", string(ioBytes), "error=", err)
			}
		}
		if io.ModifyItem.ItemInputRef == nil {
			pio.ModifyItem.ItemInputRef = -1
		} else {
			pio.ModifyItem.ItemInputRef = *io.ModifyItem.ItemInputRef
		}
		ModifyParams := GetModifyParamsFromRef(io.ModifyItem.ModifyParamsRef, t)
		pio.ModifyItem.Doubles = ModifyParams.Doubles
		pio.ModifyItem.Longs = ModifyParams.Longs
		pio.ModifyItem.Strings = ModifyParams.Strings
		wpl = append(wpl, pio)
	}

	return wpl
}

// GetModifyParamsFromRef is a function to get modifying fields from reference file
func GetModifyParamsFromRef(ref string, t *testing.T) types.ItemModifyParams {
	var iup types.ItemModifyParams
	if len(ref) == 0 {
		return iup
	}
	modBytes := ReadFile(ref, t)
	err := json.Unmarshal(modBytes, &iup)
	if err != nil {
		t.Fatal("error parsing modBytes provided via fixture Bytes=", string(modBytes), "error=", err)
	}

	return iup
}
