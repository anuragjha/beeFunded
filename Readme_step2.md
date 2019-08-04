# cs686 - Project 2: Build a private Blockchain

For this project, you are required to implement a Block structure and a BlockChain structure with some general functions. This time there's no starter code. If you are using Go, you may add an "Error" as the additional return type of all required functions.

Submission: Later we will post a github link under which you can create your own repo. And you won't have starter code for this project so you can start now. 

Resources: 
1. Python example of building a blockchain.
https://medium.com/crypto-currently/lets-build-the-tiniest-blockchain-e70965a248b (Links to an external site.)

 

## Data structures:

### 1. Block:

Each block must contain a header, and in the header there are the following fields: 
- (1) Height: int32
- (2) Timestamp: int64
The value must be in the UNIX timestamp format such as 1550013938
- (3) Hash: string.

Block’s hash is the SHA3-256 encoded value of this string(note that you have to follow this specific order): 

hash_str := string(b.Header.Height) + string(b.Header.Timestamp) + b.Header.ParentHash + b.Value.Root + string(b.Header.Size)

- (4) ParentHash: string
- (5) Size: int32
The size is the length of the byte array of the block value

Each block must have a value, which is a Merkle Patricia Trie. All the data are inserted in the MPT and then a block contains that MPT as the value. So the field definition is this: 
Value: mpt MerklePatriciaTrie

Here's the summary of block structure: 
Block{Header{Height, Timestamp, Hash, ParentHash, Size}, Value}

#### Required functions: 
If arguments or return type is not specified, feel free to define them yourself. You may change the function's name, but make a comment to indicate which function you are implementing.

- (1) Initial()
Description: This function takes arguments(such as height, parentHash, and value of MPT type) and forms a block. This is a method of the block struct.

- (2) DecodeFromJson(jsonString)
Description: This function takes a string that represents the JSON value of a block as an input, and decodes the input string back to a block instance. Note that you have to reconstruct an MPT from the JSON string, and use that MPT as the block's value. 
Argument: a string of JSON format
Return value: a block instance

- (3) EncodeToJSON()
Description: This function encodes a block instance into a JSON format string. Note that the block's value is an MPT, and you have to record all of the (key, value) pairs that have been inserted into the MPT in your JSON string. There's an example with details on Piazza. Here's a website that can encode and decode JSON string: Link (Links to an external site.)
Argument: a block or you may define this as a method of the block struct
Return value: a string of JSON format

Example of a block's JSON(decoded from JSON string):

{
    "hash":"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48",
    "timeStamp":1234567890,
    "height":1,
    "parentHash":"genesis",
    "size":1174,
    "mpt":{
        "charles":"ge",
        "hello":"world"
    }
}

### 2. BlockChain:

Each blockchain must contain two fields described below. Don't change the name or the data type. 
- (1) Chain: map[int32][]Block
This is a map which maps a block height to a list of blocks. The value is a list so that it can handle the forks.
- (2) Length: int32
Length equals to the highest block height.

#### Required functions:
If arguments or return type is not specified, feel free to define them yourself. You may change the function's name, but make a comment to indicate which function you are implementing.

- (1) Get(height)
Description: This function takes a height as the argument, returns the list of blocks stored in that height or None if the height doesn't exist.
Argument: int32
Return type: []Block

- (2) Insert(block)
Description: This function takes a block as the argument, use its height to find the corresponding list in blockchain's Chain map. If the list has already contained that block's hash, ignore it because we don't store duplicate blocks; if not, insert the block into the list. 
Argument: block

- (3) EncodeToJSON(self)
Description: This function iterates over all the blocks, generate blocks' JsonString by the function you implemented previously, and return the list of those JsonStritgns. 
Return type: string

Example of a blockchain's JSON:
```
[
    {
        "hash":"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48",
        "timeStamp":1234567890,
        "height":1,
        "parentHash":"genesis",
        "size":1174,
        "mpt":{
            "hello":"world",
            "charles":"ge"
        }
    },
    {
        "hash":"24cf2c336f02ccd526a03683b522bfca8c3c19aed8a1bed1bbc23c33cd8d1159",
        "timeStamp":1234567890,
        "height":2,
        "parentHash":"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48",
        "size":1231,
        "mpt":{
            "hello":"world",
            "charles":"ge"
        }
    }
]
```
- (4) DecodeFromJSON(self, jsonString)
Description: This function is called upon a blockchain instance. It takes a blockchain JSON string as input, decodes the JSON string back to a list of block JSON strings, decodes each block JSON string back to a block instance, and inserts every block into the blockchain. 
Argument: self, string
