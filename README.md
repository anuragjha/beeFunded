# beeFunded
#### Nakamoto Consensus based Blockchain platform running crowd funding application

#### Bockchain Specifications
Blockchain platform based on Nakamoto Consensus. Used gossip protocol for communication. Merkle patricia trie to store block data. 
###### step 1 - package(data_structure) - Data Structure for Block Data - Merkle Patricia Trie
###### step 2 - package(block and blockchain) - Defining block and chain mechanism for single system
###### step 2.5 - package(identity) - Private and Public identity with signature and other utilities
###### step 3 - package(uri_routing) - router, logger and handlers
###### step 3.5 - package(gossip_protocol) - Communication between miners - Gossip Protocol
###### step 3.7 - package(sync_blockchain) - thread safe blockchain
###### step 4 - package(pow) - Consensus used in Blockchain - POW + Nakamoto
###### step 5 - package(token) - includes tokens, wallet, transaction, borrow transaction and tx pool
###### step 5.5 - package(client) - client and blockchainholders for clients


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


#### Todo - 
- [X] Adding mechanism for fees for block producer
- [ ] Self Adjusting Difficulty - in progress
> Algo for now - Receieved a Block for a height h (put logic before checking the difficulty POW)

> Also, Algo to be put for creating a new block - (put this logic in nonce find)

 UpdateDifficulty() {
 -  find the latest n blocks (n=5 for now)
 -  find the avgTimeDiff between these 5 blocks
 -  if the avgTimeDiff > 6 sec then  decrease the difficulty level
 -  if the avgTimeDiff < 4 sec then increase the difficulty level
 
 }
 
>  Will also have to be in Download of New Miner (to get the latest Difficulty) - after that can be changed by node independently

### Assumptions - Rules under which the application works correctly
1. First Miner runs on port 6686
2. The Miner at 6686 first will create first 2 blocks only then we can start other nodes - miner as well as clients
(1st block is a init block, and 2nd block c0ntains a transaction of 10000 tokens that every miner gets when they start up)


## How to run the application


Miner and Client both uses the same initial command, with cmdline params as port number on which it wants to run on.

Miners use GET /start to begin mining process
1. miner new - ask for txPool
2. miner old - real time communicate new transaction

Clients
1. Clients use GET / to initiate client functionality
2. Client use /signup to get a key-value pair
3. Client use GET /cidpage to enter key-pair
3. After then during login use the same key pair to prove authenticity
4. Once logged in clients can initiate a ask transaction or a promise transaction, can also send tokens directly to peer.


### How to initilize a Clients
> go run main.go 7000

1. localhost:7000/ -> Enter username and submit
2. Click the signup button -> POST/signup is called
3. STORE the key-pair generated

4. On a different tab => localhost:7000/cidpage
5. Paste the stored key-pair -> and click submit button
ClientId is now stored in Client at 7000

6. Back to previous tab and back to get -> localhost:7000/
7. Click login button to get -> localhost:7000/login
8. Enter a phrase and Paste the key-pair then click the login button
9. Then use the interface to interact with other clients



### Ways to interact with interface ->


#### A. Create a Borrow requirement
1. To create a - borrow transaction
2. Enter a value in -> Amount and Fees
3. And click submit

#### B. Promise sum of amount for some Requirement
1. Decide the sum of amount to give to the Client
2. Enter that sum in -> Amount field
3. Enter Fees in -> Fees field
4. Find the TransactionId and PublicId of the Borrow transaction in bottom section -> Requirements Contact Details
5. Paste the values in -> TransactionId and -> To field
6. click the submit button

#### C. Client have enough Promises for sum of Amount against the Borrow requirement
1. Balance will be updated for all the Promisers and the Borrower, promises made cannot be reversed.

#### D. Client does not have enough Promises for sum of Amount for the Borrow requirement
1. Borrower can see the promises made against the Asked Amount
2. Promiser can request the all/partial promise back if the Promised sum has not exceeded the Borrow requirement (Added after deadline)

#### E. Transfer sum of Amount directly to another client
1. Client can enter that sum in -> Amount field
2. Enter Fees in -> Fees field
3. Enter the PublicId of the other client -> To field
4. Click on the submit button -> Balance at Sender and Receiever will be updated


## API overview
There are two actors defined for the system, Miners and Clients.

### API's for Miners
1. GET   /start 
2. GET   /show
3. GET   /upload
4. GET   /block/{height}/{hash}
5. POST  /heartbeat/receive
6. GET   /Canonical
7. GET   /showBlock/{height}
8. GET   /showBlockMpt/{height}
9. GET   /showBalanceBook
10. GET   /showTransactionPool
11. GET  /clientsignup
12. GET  /clientlogin
13. POST /txbeat/receive
14. GET  /txbeat/allprev

### Endpoints for Clients
1. GET  /
2. POST /signup
3. POST /login
4. GET  /cidpage
5. POST /setcid 
5. POST /transactionform
6. GET  /GetMyId
7. GET  /showWallet
8. GET  /bcholders
9. GET  /showbcholders


# Majors features of the Application

## How Balanced Book is built
BalanceBook is a struct which contains Book (a Key- value to store hash(PublicId) and Value(account Balance)).
1. Then we get a canonical chain and iterate over each block starting from the initial height.
2. Inside each block, iterate over transaction stored in mpt. 
3. Based on the transactions start building the book and also build Promised. 
Book is a key - value store of  hash(PublicKey) - Balance sum of tokens.
Promised is a key - value store of TransactionId - Borrowing Transaction. \
Borrowing Transaction is structutre that contains
 a. Initial Requirement Transaction
 b. Array of Promise Transaction
 c. Variable that keeps record of the total promised value.

## (1) Crypto 
Achieving data authenticity and integrity by use of Signature.

## security.go
Implements the public-private keys and Hash generation. Also creating signature and verifying signature.
Contains two data structs - Identity and PublicIdentity.
GetPublicIdentity stuct in security.go - for Node to be able to get its PublicId in new variable.

Using golang rsa PKCS1v15 to sign and verify.

1. new data structs in security.go
  - Identity {
      privateKey,
      PublicKey,
      Label,
      }
  
  - PublicIdentity {
      PublicKey,
      Label,
      }

## peerlist.go
 ### Todo - merge logic of processing PeerMap and PeerMapPid 
 PeerList - new variables in struct
  secureId    - added to contain secureId of dataType Identity.
  peerMapPid  - now also contains PeerMapPid
  PeerMapPid  - a Map of - Addr of Node (key) - and PublicIdentity (value) of peers.
  
  Methods added - 
  InjectPeerMapPidJson method - to inject receieved pidJson in receiever map
  And other methods relating PeerMapPid, parallel in logic with PeerMap
 
## heartbeat.go
 HeartBeat changed to now also include
   Pid
   SignForBlockJson - for BlockJson
   PeerMapPidJson  

## handlers.go
(Todo - sender have to encrypt the heartbeat with public key of receving peers, Receiever have to decrypt the heartbeat with private key of itself. - only when confidentiality is needed)

1) In StartHeartBeat func - AND - In SendBlockBeat func -
Add new params to PrepareHeartBeatData func call
    add signature for blockjson
    add Pid of Sender
    
(Add Encrypt heartbeat with public key of "to whom the heartbeat is being sent". (in for loop))
((2)In HeartBeatReceive func -)
(Add functionality to Decrypt heartbeat as soon as receieved)



## API
GET /uploadpids - to download the public key map from Register server to be used first time along with download blockchain.

HeartBeat send and Receieve now additionally deals with Signature of sender and its verification by receiever.

(Have to add functionality to send encrypted blockjson and decrypt blockjson - in handler.go (funcs available in security.go))


## (2) Currency
Adding mechanism to enable comodity(digital token) to be exchanged between peers.

Peers have 1000 in begining by default for now.
They can create tx, the tx then goes to tx pool.
Txs are picked by a peer from Tx pool.
Tx remain in pool until it is part of canonical chain

##### More on transactions
There are 2 type of transaction
- Borrowing Tx
- Lending Tx

Algorithm
For every Borrowing TX 
  Create a Promised Struct {  PromiseMPT<Lender TX, Lending Amount> } 
And add it in PromiseList struct { Map < BorrowingTX, Promised Struct > }

if Borrowing Tx 
  Create a Promised Struct {  PromiseMPT<Lender TX, Lending Amount> } 
  And add it in PromiseList { Map < BorrowingTX, Promised Struct > }
else if Lending Tx
  loop over PromiseList to find the matching BorrowingTXId THEN
    >if lender id is new 
      Add lender id, lended amount on Promised MPT for that Borrowing Tx in the PromiseList Map
    else if lender id is old
      Get that lender id, increase the existing amout by this new lended tokens
    >if Total Promised Amt meets the Borrowing Requirement
      Process the entry in PromiseList Map for that Borrowing TX !! processPromises(Borrowing Tx) !!(1)
        
!!(1)!!
processPromises(Borrowing Tx)
  get Entry from PromiseList struct { Map < BorrowingTX, Promised Struct > } for the corresponding Borrowing Tx
  Add the total of lended tokens 
  Remove the entry from PromiseList Map
  
---------------------------------------------------------
------------------------- some ALGO ---------------------
When Node receieves heartbeat :
  ...
  gets mpt 



---------------------------------------------------------
## transaction.go
A transaction is considered valid if tokens needed(including fees) >= balance in Book - Promised

Data structs in transaction.go
  1) Transaction<type of Tx - Borrowing or Lending ?> {
    Id - is hash of tx <<<<<<<<<
    From
    To
    AmountOfTokens <<<<<<<<<<<<<
    Timestamp
  }

  2) TransactionPool {
     list of transaction
  }

  3) TransactionBeat {
      Transaction
      FromPid
      TxSignature  <<<<<<<<
  }

Funcs ->
CreateTransaction func -> takes params From public Id, To public Id, Tokens and Timestamp and -> returns Tx.
NewTransactionBeat func -> takes params Tx, From public Id and FromSig and -> returns TransactionBeat.
CreateTransactionBeat func ->  takes params Tx and Identity and -> returns. 
AddToTransactionPool func -> takes a transaction and adds it to TransactionPool.
DeleteFromTransactionPool func -> takes transaction id and deletes it from TransactionPool.
ReadFromTransactionPool func -> takes in no. of tx to read and returns txmap of as many txs.




## wallet.go
Data structs in wallet.go
  - Wallet {
    Balance, - map[string]float64
    mutex,
  }
  where Balance is map of unit and currency amt.
  
## balanceBook.go
Data structs in balanceBook.go
  - BalanceBook {
    Book,     - mpt
    Promised, - mpt
    mutex,
  }


Funcs ->
UpdateBalanceBook
GetBalance

//todo 
a. IsBalanceEnough finc
IsBalanceEnough() {
 - takes in key and needed balance - and returns true of false based on (Book and Promised)
}
b.  generate balancebook and Promise book for a Chain
GenerateBalanceAndPromise(SBC SyncBlockchain) {
  - use function in canonical chain to get - the canonical blockchain then
  - start reading from 1st height block and read all the transactions to build up balancebook and promise book
}
c. Reading all transaction of one block 
 - get a Block 
 - convert block to key value pairs
 - for every key value pairs - check and update ... balancebook and promise

todo //

## handlers.go
Add data structs to Keep Balance, 

Need encryption of Holding account ????
if for every borrowing tx - a holding account is created.
And Lending amt is kept there until it is ready for use by borrowers.

## API

GET /ShowWallet

GET /ShowBalanceBook

GET /ShowTransactionPool

POST /transaction
  req body should contain a Transaction

POST /transactionBeatRecv
  req body should contain TransactionBeat( {Tx, FromPid, TxSignature, Hops} ) 


=====================================Client=============================================

starts normal but do not call the start api
Instead calls GET /cover api - which shows the starting page for client - func startClient()
then client can signup to get a key pair 
then store it.
After that client cal login by pasting the whole thing as -key- and -addr- and then send to any miner (for now 6686)
The miner will verify if it - if verified send the peerlist to the client. - and show Client Page (where client interacts)

data structres a Client will have ->
  - ClientId
  - SELF_ADDR from Init function
  - use here Peers of type PeerList

============================
two miners - tx beat  - 2 conditions
1. miner new - ask for txPool
2. miner old - real time communicate new transaction

when 2 miners exist - client will have two miners in bcHolders

------
keyValuePairs := p1.MerklePatriciaTrie(b.Block(uBlock).Value).GetAllKeyValuePairs()
------

When somebody signup - miner will initiate a default tx of 1000

Added functionality to handle - when sender sends more than they have
Add functionality to find promised amount for any ask req



