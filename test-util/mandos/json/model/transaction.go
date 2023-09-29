package mandosjsonmodel

// TransactionType describes the type of simulate transaction
type TransactionType int

const (
	// ScDeploy describes a transaction that deploys a new contract
	ScDeploy TransactionType = iota

	// ScCall describes a regular smart contract call
	ScCall

	// Transfer is an MOA transfer transaction without calling a smart contract
	Transfer

	// ValidatorReward is when the protocol sends a validator reward to the target account.
	// It increases the balance, but also increments "DME_Reward" in storage.
	ValidatorReward
)

// HasSender is a helper function to indicate if transaction has `to` field.
func (tt TransactionType) HasSender() bool {
	return tt != ValidatorReward
}

// HasReceiver is a helper function to indicate if transaction has receiver.
func (tt TransactionType) HasReceiver() bool {
	return tt != ScDeploy
}

// IsSmartContractTx indicates whether tx type allows an `expect` field.
func (tt TransactionType) IsSmartContractTx() bool {
	return tt == ScDeploy || tt == ScCall
}

// Transaction is a json object representing a transaction.
type Transaction struct {
	Type      TransactionType
	Nonce     JSONUint64
	Value     JSONBigInt
	From      JSONBytes
	To        JSONBytes
	Function  string
	Code      JSONBytes
	Arguments []JSONBytes
	GasPrice  JSONUint64
	GasLimit  JSONUint64
}

// TransactionResult is a json object representing an expected transaction result.
type TransactionResult struct {
	Out        []JSONCheckBytes
	Status     JSONBigInt
	Message    string
	Gas        JSONCheckUint64
	Refund     JSONCheckBigInt
	IgnoreLogs bool
	LogHash    string
	Logs       []*LogEntry
}

// LogEntry is a json object representing an expected transaction result log entry.
type LogEntry struct {
	Address    JSONBytes
	Identifier JSONBytes
	Topics     []JSONBytes
	Data       JSONBytes
}
