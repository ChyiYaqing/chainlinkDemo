package test

import (
	"testing"

	gm "github.com/bcds/go-crypto-gm"
	"github.com/bcds/gosdk/account"
	"github.com/bcds/gosdk/common"
	rpc "github.com/bcds/gosdk/rpc"
	"github.com/stretchr/testify/assert"
)

const (
	abiContract = `[{"constant":false,"inputs":[{"name":"num1","type":"uint32"},{"name":"num2","type":"uint32"}],"name":"add","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"archiveSum","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"getSum","outputs":[{"name":"","type":"uint32"}],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"increment","outputs":[],"payable":false,"type":"function"}]`
	binContract = "0x60606040526000805463ffffffff19169055341561001957fe5b5b61012a806100296000396000f300606060405263ffffffff60e060020a6000350416633ad14af38114603e57806348fe842114605c578063569c5f6d14606b578063d09de08a146091575bfe5b3415604557fe5b605a63ffffffff6004358116906024351660a0565b005b3415606357fe5b605a60c2565b005b3415607257fe5b607860d2565b6040805163ffffffff9092168252519081900360200190f35b3415609857fe5b605a60df565b005b6000805463ffffffff808216850184011663ffffffff199091161790555b5050565b6000805463ffffffff191690555b565b60005463ffffffff165b90565b6000805463ffffffff8082166001011663ffffffff199091161790555b5600a165627a7a72305820caa934a33fe993d03f87bdf39706fada68ddde78182e0110fd43e8c323d5984a0029"
)

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
	transaction := rpc.NewTransaction(guomiKey.GetAddress().Hex()).Deploy(binContract)
	transaction.Sign(guomiKey)
	ans, err := tg.DeployContract(transaction)
	if err != nil {
		t.Error(err)
	}
	t.Log("contract address: ", ans)
}
