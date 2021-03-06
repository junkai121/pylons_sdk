package fixturetest

import (
	testing "github.com/Pylons-tech/pylons_sdk/cmd/fixtures_test/evtesting"
)

// ActFunc describes the type of function used for action running test
type ActFunc func(FixtureStep, *testing.T)

var actFuncs = make(map[string]ActFunc)

// RegisterActionRunner registers action runner function
func RegisterActionRunner(action string, fn ActFunc) {
	actFuncs[action] = fn
}

// GetActionRunner get registered action runner function
func GetActionRunner(action string) ActFunc {
	return actFuncs[action]
}

// RunActionRunner execute registered action runner function
func RunActionRunner(action string, step FixtureStep, t *testing.T) {
	fn := GetActionRunner(action)
	if fn == nil {
		t.Fatalf("step with unrecognizable action found %s", step.Action)
	}
	t.MustTrue(fn != nil)
	fn(step, t)
}

// RegisterDefaultActionRunners register default test functions
func RegisterDefaultActionRunners() {
	RegisterActionRunner("fiat_item", RunFiatItem)
	RegisterActionRunner("update_item_string", RunUpdateItemString)
	RegisterActionRunner("create_cookbook", RunCreateCookbook)
	RegisterActionRunner("create_recipe", RunCreateRecipe)
	RegisterActionRunner("execute_recipe", RunExecuteRecipe)
	RegisterActionRunner("check_execution", RunCheckExecution)
	RegisterActionRunner("create_trade", RunCreateTrade)
	RegisterActionRunner("fulfill_trade", RunFulfillTrade)
	RegisterActionRunner("disable_trade", RunDisableTrade)
	RegisterActionRunner("multi_msg_tx", RunMultiMsgTx)
}
