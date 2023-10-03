## NFT Authenticator

**What is the NFT Authenticator?**

The NFT Authenticator is a type of Authenticator designed to streamline transactions by leveraging Non-Fungible Tokens (NFTs). By using the NFT Authenticator, an NFT owner can send an NFT to an address, allowing that address to transact on behalf of the original owner address. 

**Workflow with the Osmosis Token Factory**

Here's a simplified example of how the NFT Authenticator works in conjunction with the Osmosis Token Factory:

```go
1. *Alice* mints an NFT using the Osmsois Token Factory.

2. *Alice* now possesses the NFT.

3. *Alice* decides to send the NFT to *Bob*.

4. *Bob* wants to perform a transaction on behalf of *Alice*, so they sign a transaction with *Alice* as the sender (since *Alice* is the minter).

5. The NFT Authenticator checks that *Bob* has signed the transaction and the NFT in their account.

6. The transaction is successful because *Bob* has the NFT and the signature is correct.

7. *Chris* feel left out and tries to send transactions on behalf of *Alice* but fails.

8. The spend limit authenticator and NFT authenticators together can be used with an AllOf authenticator to ensure secure and controlled transactions.
```

This workshop serves as an introduction to the concepts of Account Abstraction and the NFT Authenticator. Feel free to explore the [Osmosis Authenticator Module](https://github.com/osmosis-labs/osmosis/blob/account-abstraction-main/x/authenticator) repository for more in-depth examples, and practical applications of these concepts.
  
