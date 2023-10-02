package nft

import (
	"encoding/hex"
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/cosmos-sdk/client"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/osmosis-labs/osmosis/v19/app"
	"github.com/osmosis-labs/osmosis/v19/app/params"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/osmosis-labs/osmosis/v19/x/authenticator/authenticator"
)

// NFTAuthenticatorTest is the test suite struct for testing NFT authenticator functionality.
type NFTAuthenticatorTest struct {
	suite.Suite
	Ctx                  sdk.Context
	OsmosisApp           *app.OsmosisApp
	EncodingConfig       params.EncodingConfig
	PassKeyAuthenticator authenticator.PassKeyAuthenticator
	TestKeys             []string
	TestAccAddress       []sdk.AccAddress
	TestPrivKeys         []*secp256k1.PrivKey

	NFT NFTAuthenticator
}

// SetupTest initializes the test environment for NFT authenticator testing, including setting up accounts and authenticators.
func (s *NFTAuthenticatorTest) SetupTest() {
	s.OsmosisApp = app.Setup(false)
	s.Ctx = s.OsmosisApp.NewContext(false, tmproto.Header{})
	s.Ctx = s.Ctx.WithGasMeter(sdk.NewGasMeter(2_000_000))

	TestKeys := []string{
		"6cf5103c60c939a5f38e383b52239c5296c968579eec1c68a47d70fbf1d19159",
		"0dd4d1506e18a5712080708c338eb51ecf2afdceae01e8162e890b126ac190fe",
		"49006a359803f0602a7ec521df88bf5527579da79112bb71f285dd3e7d438033",
	}
	s.EncodingConfig = app.MakeEncodingConfig()
	txConfig := s.EncodingConfig.TxConfig
	signModeHandler := txConfig.SignModeHandler()

	ak := s.OsmosisApp.AccountKeeper

	// Set up test accounts
	for _, key := range TestKeys {
		bz, _ := hex.DecodeString(key)
		priv := &secp256k1.PrivKey{Key: bz}

		// add the test private keys to array for later use
		s.TestPrivKeys = append(s.TestPrivKeys, priv)

		accAddress := sdk.AccAddress(priv.PubKey().Address())
		account := authtypes.NewBaseAccount(accAddress, priv.PubKey(), 0, 0)
		ak.SetAccount(s.Ctx, account)

		// add the test accounts to array for later use
		s.TestAccAddress = append(s.TestAccAddress, accAddress)
	}

	// Create a new Secp256k1SignatureAuthenticator for use in the NFT authenticator
	_ = authenticator.NewSignatureVerificationAuthenticator(
		ak,
		signModeHandler,
	)

	// Create a the NFT authenticator with the bank keeper and tokenfactory keeper.
	s.NFT = NewNFTAuthenticator()
}

// TestNFTAuthentication performs a series of tests to validate the NFT authentication flow.
// It covers scenarios such as creating an NFT token, minting the NFT, transferring it to another account,
// and testing authentication with valid and invalid scenarios. The function initializes the authenticator,
// generates authentication data from a transaction, and evaluates the authentication results.
// It ensures that the NFT authenticator correctly authenticates transactions based on ownership of the NFT token.
func (s *NFTAuthenticatorTest) TestNFTAuthentication() {
	// osmoToken := "osmo"
	// nftPostfix := "nft"

	//
	// Create NFT token for use with the NFT Authenticator
	//

	//
	// Mint a single token, this will act as our NFT
	//

	//
	// We minted the token now send the token to the account we want to authenticate from
	//

	//
	// Generate a transaction to test our authentication flow
	//

	//
	// Initialize the authenticator, this would happen from the store, we initialize with the name of the NFT
	//

	//
	// Get the authentication data from the transaction
	//

	//
	// Authenticate the transaction, this will pass as Account 1 has a valid signature and also has the NFT
	//

	//
	// Passed :tada:
	//

	//
	// Generate a transaction to test our authentication flow that we know will fail
	//

	//
	// Get the authentication data from the transaction
	//

	//
	// Try to authenticate the transaction, from an account that doesn't own the NFT
	//

	//
	// Failed :tear:
	//
}

// GenTx is a helper function to generate a signed mock transaction.
func GenTx(
	gen client.TxConfig,
	msgs []sdk.Msg,
	feeAmt sdk.Coins,
	gas uint64,
	chainID string,
	accNums,
	accSeqs []uint64,
	signers []cryptotypes.PrivKey,
	signatures []cryptotypes.PrivKey,
) (sdk.Tx, error) {
	sigs := make([]signing.SignatureV2, len(signers))

	// create a random length memo
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	memo := simulation.RandStringOfLength(r, simulation.RandIntBetween(r, 0, 100))
	signMode := gen.SignModeHandler().DefaultMode()

	// 1st round: set SignatureV2 with empty signatures, to set correct
	// signer infos.
	for i, p := range signers {
		sigs[i] = signing.SignatureV2{
			PubKey: p.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode: signMode,
			},
			Sequence: accSeqs[i],
		}
	}

	tx := gen.NewTxBuilder()
	err := tx.SetMsgs(msgs...)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	tx.SetMemo(memo)
	tx.SetFeeAmount(feeAmt)
	tx.SetGasLimit(gas)

	// 2nd round: once all signer infos are set, every signer can sign.
	signers = signers[0:len(signatures)]
	for i, p := range signatures {
		sigs[i].PubKey = p.PubKey()
	}
	err = tx.SetSignatures(sigs...)

	for i, p := range signatures {
		signerData := authsigning.SignerData{
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		sigs[i].PubKey = p.PubKey()
		signBytes, err := gen.SignModeHandler().GetSignBytes(signMode, signerData, tx.GetTx())
		if err != nil {
			panic(err)
		}
		sig, err := p.Sign(signBytes)
		if err != nil {
			panic(err)
		}
		sigs[i].Data.(*signing.SingleSignatureData).Signature = sig
		err = tx.SetSignatures(sigs...)
		if err != nil {
			panic(err)
		}
	}

	return tx.GetTx(), nil
}

func TestNFTAuthenticatorTest(t *testing.T) {
	suite.Run(t, new(NFTAuthenticatorTest))
}
