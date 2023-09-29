package mandosjsonwrite

import (
	mj "github.com/kalyan3104/dme-vm-util/test-util/mandos/json/model"
	oj "github.com/kalyan3104/dme-vm-util/test-util/orderedjson"
)

// ScenarioToJSONString converts a scenario object to its JSON representation.
func ScenarioToJSONString(scenario *mj.Scenario) string {
	jobj := ScenarioToOrderedJSON(scenario)
	return oj.JSONString(jobj)
}

// ScenarioToOrderedJSON converts a scenario object to an ordered JSON object.
func ScenarioToOrderedJSON(scenario *mj.Scenario) oj.OJsonObject {
	scenarioOJ := oj.NewMap()

	if len(scenario.Name) > 0 {
		scenarioOJ.Put("name", stringToOJ(scenario.Name))
	}

	if len(scenario.Comment) > 0 {
		scenarioOJ.Put("comment", stringToOJ(scenario.Comment))
	}

	if !scenario.CheckGas {
		ojFalse := oj.OJsonBool(false)
		scenarioOJ.Put("checkGas", &ojFalse)
	}

	var stepOJList []oj.OJsonObject

	for _, generalStep := range scenario.Steps {
		stepOJ := oj.NewMap()
		stepOJ.Put("step", stringToOJ(generalStep.StepTypeName()))
		switch step := generalStep.(type) {
		case *mj.ExternalStepsStep:
			stepOJ.Put("path", stringToOJ(step.Path))
		case *mj.SetStateStep:
			if len(step.Comment) > 0 {
				stepOJ.Put("comment", stringToOJ(step.Comment))
			}
			if len(step.Accounts) > 0 {
				stepOJ.Put("accounts", accountsToOJ(step.Accounts))
			}
			if len(step.NewAddressMocks) > 0 {
				stepOJ.Put("newAddresses", newAddressMocksToOJ(step.NewAddressMocks))
			}
			if step.PreviousBlockInfo != nil {
				stepOJ.Put("previousBlockInfo", blockInfoToOJ(step.PreviousBlockInfo))
			}
			if step.CurrentBlockInfo != nil {
				stepOJ.Put("currentBlockInfo", blockInfoToOJ(step.CurrentBlockInfo))
			}
			if len(step.BlockHashes) > 0 {
				stepOJ.Put("blockHashes", blockHashesToOJ(step.BlockHashes))
			}
		case *mj.CheckStateStep:
			if len(step.Comment) > 0 {
				stepOJ.Put("comment", stringToOJ(step.Comment))
			}
			stepOJ.Put("accounts", checkAccountsToOJ(step.CheckAccounts))
		case *mj.TxStep:
			if len(step.TxIdent) > 0 {
				stepOJ.Put("txId", stringToOJ(step.TxIdent))
			}
			if len(step.Comment) > 0 {
				stepOJ.Put("comment", stringToOJ(step.Comment))
			}
			stepOJ.Put("tx", transactionToScenarioOJ(step.Tx))
			if step.Tx.Type.IsSmartContractTx() && step.ExpectedResult != nil {
				stepOJ.Put("expect", resultToOJ(step.ExpectedResult))
			}
		}

		stepOJList = append(stepOJList, stepOJ)
	}

	stepsOJ := oj.OJsonList(stepOJList)
	scenarioOJ.Put("steps", &stepsOJ)

	return scenarioOJ
}

func transactionToScenarioOJ(tx *mj.Transaction) oj.OJsonObject {
	transactionOJ := oj.NewMap()
	if tx.Type.HasSender() {
		transactionOJ.Put("from", byteArrayToOJ(tx.From))
	}
	if tx.Type.HasReceiver() {
		transactionOJ.Put("to", byteArrayToOJ(tx.To))
	}
	transactionOJ.Put("value", bigIntToOJ(tx.Value))
	if tx.Type == mj.ScCall {
		transactionOJ.Put("function", stringToOJ(tx.Function))
	}
	if tx.Type == mj.ScDeploy {
		transactionOJ.Put("contractCode", byteArrayToOJ(tx.Code))
	}

	if tx.Type == mj.ScCall || tx.Type == mj.ScDeploy {
		var argList []oj.OJsonObject
		for _, arg := range tx.Arguments {
			argList = append(argList, byteArrayToOJ(arg))
		}
		argOJ := oj.OJsonList(argList)
		transactionOJ.Put("arguments", &argOJ)
	}

	if tx.Type.IsSmartContractTx() {
		transactionOJ.Put("gasLimit", uint64ToOJ(tx.GasLimit))
		transactionOJ.Put("gasPrice", uint64ToOJ(tx.GasPrice))
	}

	return transactionOJ
}

func newAddressMocksToOJ(newAddressMocks []*mj.NewAddressMock) oj.OJsonObject {
	var namList []oj.OJsonObject
	for _, namEntry := range newAddressMocks {
		namOJ := oj.NewMap()
		namOJ.Put("creatorAddress", byteArrayToOJ(namEntry.CreatorAddress))
		namOJ.Put("creatorNonce", uint64ToOJ(namEntry.CreatorNonce))
		namOJ.Put("newAddress", byteArrayToOJ(namEntry.NewAddress))
		namList = append(namList, namOJ)
	}
	namOJList := oj.OJsonList(namList)
	return &namOJList
}

func blockInfoToOJ(blockInfo *mj.BlockInfo) oj.OJsonObject {
	blockInfoOJ := oj.NewMap()
	if len(blockInfo.BlockTimestamp.Original) > 0 {
		blockInfoOJ.Put("blockTimestamp", uint64ToOJ(blockInfo.BlockTimestamp))
	}
	if len(blockInfo.BlockNonce.Original) > 0 {
		blockInfoOJ.Put("blockNonce", uint64ToOJ(blockInfo.BlockNonce))
	}
	if len(blockInfo.BlockRound.Original) > 0 {
		blockInfoOJ.Put("blockRound", uint64ToOJ(blockInfo.BlockRound))
	}
	if len(blockInfo.BlockEpoch.Original) > 0 {
		blockInfoOJ.Put("blockEpoch", uint64ToOJ(blockInfo.BlockEpoch))
	}

	return blockInfoOJ
}
