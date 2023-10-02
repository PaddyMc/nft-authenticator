package nft

import (
	"fmt"
	"testing"

	"github.com/osmosis-labs/osmosis/v19/app"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ibctesting "github.com/cosmos/ibc-go/v4/testing"

	"github.com/osmosis-labs/osmosis/v19/app/apptesting"
	"github.com/osmosis-labs/osmosis/v19/tests/osmosisibctesting"

	tokenfactorytypes "github.com/osmosis-labs/osmosis/v19/x/tokenfactory/types"

	"github.com/stretchr/testify/suite"
)

// AuthenticatorSuite is the test suite struct for integration tests related to authenticator functionality.
type AuthenticatorSuite struct {
	apptesting.KeeperTestHelper

	coordinator *ibctesting.Coordinator

	chainA *osmosisibctesting.TestChain
	app    *app.OsmosisApp

	PrivKeys []cryptotypes.PrivKey
	Account  authtypes.AccountI
}

// SetupTest initializes the test environment for integration testing of authenticator functionality.
func (s *AuthenticatorSuite) SetupTest() {
	// Use the osmosis custom function for creating an osmosis app
	ibctesting.DefaultTestingAppInit = osmosisibctesting.SetupTestingApp

	// Here we create the app using ibctesting
	s.coordinator = ibctesting.NewCoordinator(s.T(), 1)
	s.chainA = &osmosisibctesting.TestChain{
		TestChain: s.coordinator.GetChain(ibctesting.GetChainID(1)),
	}
	s.app = s.chainA.GetOsmosisApp()

	// Initialize three private keys for testing
	s.PrivKeys = make([]cryptotypes.PrivKey, 3)
	for i := 0; i < 3; i++ {
		s.PrivKeys[i] = secp256k1.GenPrivKey()
	}

	// Initialize a test account with the first private key
	s.Account = s.CreateAccount(s.PrivKeys[0], 500_000)

	// Reduce the gas costs for creating a token in the token factory
	s.app.TokenFactoryKeeper.SetParams(
		s.chainA.GetContext(),
		tokenfactorytypes.NewParams(sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(1))), 1_000),
	)
}

// TestRoundTripAliceBobAndChris tests the full round-trip of using NFT authenticators with three different users:
// Alice, Bob, and Chris. It covers scenarios such as allowing Bob to transact on behalf of Alice after Alice sends him the NFT,
// adding authenticators to Alice's account, creating NFTs, minting NFTs, sending coins, and authenticating transactions.
func (s *AuthenticatorSuite) TestRoundTripAliceBobAndChris() {
	//	osmoToken := "osmo"
	//	nftPostfix := "nft"

	//
	// Alice is the original minter of the NFT and will allow Bob to transact on her behalf
	// by sending the token to Bob
	//
	Alice := s.PrivKeys[0]
	AliceAcc := s.Account

	//
	// Bob will receive the NFT and be able to transact on behalf of Alice
	//
	Bob := s.PrivKeys[1]

	//
	// Chris will send transactions on behalf of Alice which will always fail
	//
	Chris := s.PrivKeys[2]

	//
	// Create a new Secp256k1SignatureAuthenticator for use in the NFT authenticator
	//

	fmt.Println(Alice)
	fmt.Println(AliceAcc)
	fmt.Println(Bob)
	fmt.Println(Chris)

	//
	// Create a the NFT authenticator with the bank keeper and tokenfactory keeper.
	//

	//
	// Register both Authenticators with the AuthenticatorManager
	//

	//
	// Add a SignatureVerificationAuthenticator to Alices account
	//

	//
	// Add a NFTAuthenticator to Alices account, here we specify the name of the NFT to use (fullNFTDenom)
	//

	//
	// Create the NFT in the tokenfactory
	//

	//
	// Create the Mint NFT message
	//

	//
	// Mint the NFT using the tokenfactory
	//

	//
	// Error sending message from Bob on behalf of Alice
	//

	//
	// Send the NFT from Alice to user Bob
	//

	//
	// Success as the NFT authenticator worked as expected for Bob
	//

	//
	// Failed as the NFT authenticator worked as expected for Chris
	//
}

// CreateAccount creates a test account with the provided private key and funds it with the specified amount of coins.
// It returns the created account.
func (s *AuthenticatorSuite) CreateAccount(privKey cryptotypes.PrivKey, amount int) authtypes.AccountI {
	accountAddr := sdk.AccAddress(privKey.PubKey().Address())
	// Fund the account
	coins := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, int64(amount)))
	err := s.app.BankKeeper.SendCoins(s.chainA.GetContext(), s.chainA.SenderAccount.GetAddress(), accountAddr, coins)
	s.Require().NoError(err, "Failed to send bank tx to account")

	return s.app.AccountKeeper.GetAccount(s.chainA.GetContext(), accountAddr)
}

// pks is an alias for a slice of private keys.
type pks = []cryptotypes.PrivKey

// TestAuthenticatorSuite runs the AuthenticatorSuite test suite.
func TestAuthenticatorSuite(t *testing.T) {
	suite.Run(t, new(AuthenticatorSuite))
}
