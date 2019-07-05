# beeFunded
peer-to-peer-lending-on-blockchain
### (usecase - decentralized crowd funding platform)

## Bockchain Specifications
#### step1 (data_structure) - Data Structure for Block Data - Merkle Patricia Trie
#### - Communication between miners - Gossip Protocol
#### - Consensus used in Blockchain - POW + Nakamoto


# Application has following features
#### Every User is both, a Borrower and a Lender.
1. Borrower can 'ask' for a sum of tokens. Each 'ask' has
2. This 'ask' request is share among network users.
3. Lender can choose a 'ask' request and 'promise' and make a promise for a sum of money.
4. All previous promises for the specific 'ask' request are tracked, and if the sum of promises mount to 'asked' tokens, then money is transferred for appropriate users.
5. Borrower can then use 'pay back' option peer to peer token transfer to pay back the lenders within a specified time.

Every 'Ask' request should have
- token required by the borrower
- interest rate (for the lenders once tokens are transferred to the borrower)
- time period within which the token along with interest will be returned (once the money is transferred to the borrower)
Available balance is different from Actual balance. Available balance takes into account the sum of tokens already promised.
Available balance is used for checking validity of transaction.
