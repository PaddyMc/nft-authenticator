package nft

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/osmosis-labs/osmosis/osmomath"
	tokenfactorykeeper "github.com/osmosis-labs/osmosis/v19/x/tokenfactory/keeper"

	"github.com/osmosis-labs/osmosis/v19/x/authenticator/authenticator"
	"github.com/osmosis-labs/osmosis/v19/x/authenticator/iface"
)

// Compile time type assertion NFTAuthenticator struct and data
var _ iface.Authenticator = &NFTAuthenticator{}
var _ iface.AuthenticatorData = &NFTAuthData{}

const (
	// NFTAuthenticatorType represents a type of authenticator
	NFTAuthenticatorType = "NFTAuthenticator"
)

// NFTAuthenticator struct contains all the necessary data to enable the
// Authenticator to verify signatures and check if a user has a NFT
type NFTAuthenticator struct {
	bankKeeper  bankkeeper.Keeper
	tokenKeeper tokenfactorykeeper.Keeper
	sva         authenticator.SignatureVerificationAuthenticator
	denom       string
}

// Type returns the NFTAuthenticatorType, this is used when an authenticator is added
// and when we have branching paths in the code we switch on the Type field of the authenticator
func (na NFTAuthenticator) Type() string {
	return NFTAuthenticatorType
}

func (na NFTAuthenticator) StaticGas() uint64 {
	// For every NFTAuthenticator we consume 1000 gas, plus signature verification
	return 1000
}

// NewNFTAuthenticator creates a new with the correct keeper needed to function
// correctly, this is added to the authentication manager when the applciation is
// started
func NewNFTAuthenticator(
	bankKeeper bankkeeper.Keeper,
	tokenKeeper tokenfactorykeeper.Keeper,
	sva authenticator.SignatureVerificationAuthenticator,
) NFTAuthenticator {
	return NFTAuthenticator{
		bankKeeper:  bankKeeper,
		tokenKeeper: tokenKeeper,
		sva:         sva,
	}
}

// Initialize is used after we get authenticator data from the store,
// we initialize data
func (na NFTAuthenticator) Initialize(
	data []byte,
) (iface.Authenticator, error) {
	na.denom = string(data)
	return na, nil
}

// NFTSignatureData is used to package all the signature data and the tx
// for use in the Authenticate function, since we wrap the SignatureVerificationAuthenticator
// we use the same signature data
type NFTAuthData = authenticator.SignatureData

// GetAuthenticationData parses the signers and signatures from a transactiom
// then returns a indexed list of both signers and signatures, we use the SignatureVerificationAuthenticator
// GetAuthenticationData function since no data is needed from the transaction for checking who owns the NFT
// NOTE: position in the array is used to associate the signer and signature
func (na NFTAuthenticator) GetAuthenticationData(
	ctx sdk.Context,
	tx sdk.Tx,
	messageIndex int,
	simulate bool,
) (iface.AuthenticatorData, error) {
	return na.sva.GetAuthenticationData(ctx, tx, messageIndex, simulate)
}

// Authenticate takes an NFTVerificationData struct and validates
// each signer and signature using signature verification, then
// ensure that the signer has the NFT that enables the use of the
// original creators account
func (na NFTAuthenticator) Authenticate(
	ctx sdk.Context,
	account sdk.AccAddress,
	msg sdk.Msg,
	authenticationData iface.AuthenticatorData,
) iface.AuthenticationResult {
	sigPubKey := authenticationData.(authenticator.SignatureData).Signatures[0].PubKey

	// Set the public key to the signature public key
	na.sva.PubKey = sigPubKey

	// Get the signer address
	signerAddress := sdk.AccAddress(sigPubKey.Address())

	// Authenticate the signature
	authenticationResult := na.sva.Authenticate(ctx, signerAddress, msg, authenticationData)
	if !authenticationResult.IsAuthenticated() {
		return authenticationResult
	}

	// Get the balances for the account of the signer
	balance := na.bankKeeper.GetBalance(ctx, signerAddress, na.denom)

	// If account doen't contain the NFT return
	if !balance.Amount.Equal(osmomath.NewInt(1)) {
		return iface.NotAuthenticated()
	}

	// Successful authentication for the NFT holder
	return authenticationResult
}

// Track is used for authenticators to track any information they may need regardless of how the transaction is
// authenticated. For instance, if a message is authenticated via authz, ICA, or similar, those entry points should
// call authenticator.Track(...) so that the authenticator can know that the account has executed a specific message.
func (na NFTAuthenticator) Track(ctx sdk.Context, account sdk.AccAddress, msg sdk.Msg) error {
	// Track any necessary information (if applicable)
	return nil
}

// OnAuthenticatorAdded is called when an authenticator is added to an account. If the data is not properly formatted
// or the authenticator is not compatible with the account, an error should be returned.
func (na NFTAuthenticator) OnAuthenticatorAdded(ctx sdk.Context, account sdk.AccAddress, data []byte) error {
	// Handle actions when an authenticator is added to an account (if needed)
	return nil
}

// OnAuthenticatorRemoved is called when an authenticator is removed from an account.
// This can be used to update any global data that the authenticator is tracking or to prevent removal
// by returning an error.
// Removal prevention should be used sparingly and only when absolutely necessary.
func (na NFTAuthenticator) OnAuthenticatorRemoved(ctx sdk.Context, account sdk.AccAddress, data []byte) error {
	// Handle actions when an authenticator is removed from an account (if needed)
	return nil
}

// ConfirmExecution is employed in the post-handler function to enforce transaction rules,
// such as spending and transaction limits. It accesses the account's owned state to store
// and verify these values.
func (na NFTAuthenticator) ConfirmExecution(
	ctx sdk.Context,
	account sdk.AccAddress,
	msg sdk.Msg,
	authenticationData iface.AuthenticatorData,
) iface.ConfirmationResult {
	// Confirm the execution of the transaction (if needed)
	return iface.Confirm()
}
