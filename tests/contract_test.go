package test

import (
	"github.com/bcds/gosdk/abi"
	"math/big"
	"strings"
	"testing"

	gm "github.com/bcds/go-crypto-gm"
	"github.com/bcds/gosdk/account"
	"github.com/bcds/gosdk/common"
	rpc "github.com/bcds/gosdk/rpc"
	"github.com/stretchr/testify/assert"
)

const (
	// contract/1_Storage.sol
	abiStorageContract = `[{"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	binStorageContract = "608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100a1565b60405180910390f35b610073600480360381019061006e91906100ed565b61007e565b005b60008054905090565b8060008190555050565b6000819050919050565b61009b81610088565b82525050565b60006020820190506100b66000830184610092565b92915050565b600080fd5b6100ca81610088565b81146100d557600080fd5b50565b6000813590506100e7816100c1565b92915050565b600060208284031215610103576101026100bc565b5b6000610111848285016100d8565b9150509291505056fea2646970667358221220d645d93141b4f4c9539241aaf7fb78cc2db954445983dfb230bbc8fe488608ee64736f6c63430008130033"

	// contract
	abiGetChainIDContract = `[{"inputs":[],"name":"getChainID","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
	binGetChainIDContract = "608060405234801561001057600080fd5b5060ba8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063564b81ef14602d575b600080fd5b60336047565b604051603e91906061565b60405180910390f35b6000804690508091505090565b605b81607a565b82525050565b6000602082019050607460008301846054565b92915050565b600081905091905056fea2646970667358221220532c712583a56469ca2acf4d442c4b1f0cb840928f2ed0ff115a8892515d9c3464736f6c63430008000033"

	abiLuckyNumberContract = `[{"inputs":[],"stateMutability":"payable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bool","name":"isWinner","type":"bool"},{"indexed":true,"internalType":"address","name":"player","type":"address"}],"name":"LotteryEvent","type":"event"},{"inputs":[{"internalType":"uint256","name":"_number","type":"uint256"}],"name":"guessNumber","outputs":[],"stateMutability":"payable","type":"function"}]`
	binLuckyNumberContract = "60806040526103f26000819055506101c38061001c6000396000f3fe60806040526004361061001e5760003560e01c8063b438d01814610023575b600080fd5b61003d6004803603810190610038919061012a565b61003f565b005b600054811461009c573373ffffffffffffffffffffffffffffffffffffffff167f104b852b5326c1819a88a9f47e6adad4b5db2071bb33e66fa4f2a36605f35a19600060405161008f9190610172565b60405180910390a26100ec565b3373ffffffffffffffffffffffffffffffffffffffff167f104b852b5326c1819a88a9f47e6adad4b5db2071bb33e66fa4f2a36605f35a1960016040516100e39190610172565b60405180910390a25b50565b600080fd5b6000819050919050565b610107816100f4565b811461011257600080fd5b50565b600081359050610124816100fe565b92915050565b6000602082840312156101405761013f6100ef565b5b600061014e84828501610115565b91505092915050565b60008115159050919050565b61016c81610157565b82525050565b60006020820190506101876000830184610163565b9291505056fea264697066735822122069156b06a093007683ec6afc8859fedee919a03a5537a8036fcf03f1230067ce64736f6c63430008130033"
)

// Deploy contract : contract_address: 0xd0c3238722541321870a0a2c8d083278df3e3066
func TestContractDeployContract(t *testing.T) {
	g := rpc.NewGRPC()
	tg, err := g.NewContractGrpc(rpc.ClientOption{
		StreamNumber: 1,
	})
	assert.Nil(t, err)
	defer tg.Close()

	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey := &account.SM2Key{
		SM2PrivateKey: &gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}
	// default evmType : EVM
	transaction := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Deploy(binStorageContract)
	transaction.Sign(guomiKey)
	ans, err := tg.DeployContractReturnReceipt(transaction)
	if err != nil {
		t.Error(err)
	}
	t.Log("contract address: ", ans.ContractAddress)
}

func TestContractGrpc_InvokeContractReturnReceipt(t *testing.T) {
	g := rpc.NewGRPC()
	tg, err := g.NewContractGrpc(rpc.ClientOption{
		StreamNumber: 1,
	})
	assert.Nil(t, err)
	defer tg.Close()

	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey := &account.SM2Key{
		SM2PrivateKey: &gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}
	addrTransaction := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Deploy(binLuckyNumberContract)
	addrTransaction.Sign(guomiKey)
	ans, err := tg.DeployContractReturnReceipt(addrTransaction)
	if err != nil {
		t.Error(err)
	}
	ABI, err := abi.JSON(strings.NewReader(abiLuckyNumberContract))
	if err != nil {
		t.Error(err)
	}
	// getChainID
	packedRetrieve, err := ABI.Pack("guessNumber", big.NewInt(1000))
	if err != nil {
		t.Error(err)
	}
	transaction := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Invoke(ans.ContractAddress, packedRetrieve)
	transaction.Sign(guomiKey)
	ret, err := tg.InvokeContractReturnReceipt(transaction)
	if err != nil {
		t.Error(err)
	}

	t.Log(ret)
}

func TestContractGrpc_InvokeContractNoDeployReturnReceipt(t *testing.T) {
	contractAddr := "0xbf5cbde4106fe599cfa2427ebfbdcc19c999aca5"
	g := rpc.NewGRPC()
	tg, err := g.NewContractGrpc(rpc.ClientOption{
		StreamNumber: 1,
	})
	assert.Nil(t, err)
	defer tg.Close()

	guomiPri := "6153af264daa4763490f2a51c9d13417ef9f579229be2141574eb339ee9b9d2a"
	pri := new(gm.SM2PrivateKey)
	pri.FromBytes(common.FromHex(guomiPri), 0)
	guomiKey := &account.SM2Key{
		SM2PrivateKey: &gm.SM2PrivateKey{
			K:         pri.K,
			PublicKey: pri.CalculatePublicKey().PublicKey,
		},
	}
	ABI, err := abi.JSON(strings.NewReader(abiLuckyNumberContract))
	if err != nil {
		t.Error(err)
	}
	// getChainID
	packedRetrieve, err := ABI.Pack("guessNumber", big.NewInt(1010))
	if err != nil {
		t.Error(err)
	}
	transaction := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Invoke(contractAddr, packedRetrieve)
	transaction.Sign(guomiKey)
	ret, err := tg.InvokeContractReturnReceipt(transaction)
	if err != nil {
		t.Error(err)
	}

	t.Log(ret)
}
