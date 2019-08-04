# CS686-Project 1: Build a Merkle Patricia Tree

Ethereum uses a Merkle Patricia Tree (Links to an external site.) to store the transaction data in a block. By organizing the transaction data in a Merkle Patricia Tree, any block with fraudulent transactions would not match the tree's root hash. Build your own implementation of a Merkle Patricia Trie, following the specifications at the Ethereum wiki (Links to an external site.). Be mindful that you will use this code in the subsequent projects. 

 

Starter Code:

Rust: https://github.com/CharlesGe129/cs686_blockchain_P1_Rust_skeleton.git (Links to an external site.)

Go: https://github.com/CharlesGe129/cs686_blockchain_P1_Go_skeleton.git (Links to an external site.)

 

Submission

All your github repositories should be PRIVATE repositories. Before the deadline, please share your repository with TA. You may add TA's github(https://github.com/CharlesGe129) (Links to an external site.) as collaborator.

 

Helpful resources:

Rust: 

https://doc.rust-lang.org/rust-by-example/index.html (Links to an external site.)

https://doc.rust-lang.org/std/vec/?search= (Links to an external site.)

Go:

https://tour.golang.org/welcome/1 (Links to an external site.)

https://golang.org/pkg/strings/#HasPrefix (Links to an external site.)

 

### Project 1 specification
####
For this project, implement a Merkle Patricia Trie according to this Link (Links to an external site.) and instructor's lectures. The examples in this specification are provided in both Rust and Go. Skeleton code and some help functions will be provided. You should not change the skeleton code. If you have different ideas that require to change the skeleton code, see the TA *before* submission. Any modification to the skeleton code without the TA's approval will result in point deduction.
####
You have to choose between Go and Rust. Skeleton code is provided for both language and please write and pass your own tests before submitting the final version of code.

#### You need to implement five features of the Merkle Patricia Trie:

### 1. Get(key) -> value
Description: The Get function takes a key as argument, traverses down the Merkle Patricia Trie to find the value, and returns it. If the key doesn't exist, it will return an empty string.(for the Go version: if the key is nil, Get returns an empty string.)
##### Arguments: key (string)
##### Return: the value stored for that key (string).
##### Rust function definition: fn get(&mut self, key: &str) -> String
##### Go function definition: func (mpt *MerklePatriciaTrie) Get(key string) string

### 2. Insert(key, value)
Description: The Insert function takes a key and value as arguments. It will traverse  Merkle Patricia Trie, find the right place to insert the value, and do the insertion.(for the Go version: you can assume the key and value will never be nil.)
##### Arguments: key (string), value (string)
##### Return: string
##### Rust function definition: fn insert(&mut self, key: &str, new_value: &str)
##### Go function definition: func (mpt *MerklePatriciaTrie) Insert(key string, new_value string)

### 3. Delete(key)
Description: The Delete function takes a key as argument, traverses the Merkle Patricia Trie and finds that key. If the key exists, delete the corresponding value and re-balance the trie if necessary, then return an empty string; if the key doesn't exist, return "path_not_found".
##### Arguments: key (string)
##### Return: string
##### Rust function definition: fn delete(&mut self, key: &str) -> String
##### Go function definition: func (mpt *MerklePatriciaTrie) Delete(key string) string

### 4. compact_encode()
The compact_encode function takes an array of numbers as input (each number is between 0 and 15 included, representing a single hex digit), and returns an array of numbers according to the compact encoding rules in the github wiki page under "Compact encoding of hex sequence with optional terminator"). Each number in the output is between 0 and 255 included (representing an ASCII-encoded letter, or for the first value it represents the node type as per the wiki page). You may find a Python version in this Link (Links to an external site.), but be mindful that the return type is different!
##### Arguments: hex_array(array of u8)
##### Return: array of u8
##### Rust function definition: compact_encode(hex_array: Vec<u8>) -> Vec<u8>
##### Example: input=[1, 6, 1], encoded_array=[1, 1, 6, 1], output=[17, 97]

### 5. compact_decode()
Description: This function reverses the compact_encode() function. 
##### Arguments: hex_array(array of u8)
##### Return: array of u8
##### Rust function definition: compact_decode(encoded_arr: Vec<u8>) -> Vec<u8>
##### Example: input=[17, 97], output=[1, 6, 1]

#### Other help functions:

##### 1. fn hash_node(node: &Node) -> String
Description: This function takes a node as the input, hash the node and return the hashed string.

If you use Golang, please follow this link to install the SHA3-256 package: https://github.com/golang/crypto (Links to an external site.)

If you use Rust, the package dependency is written into Cargo.toml, so directly build and run the project should be fine. 

Classes specification: 
In this project, there are two pre-defined classes. Both classes are defined in skeleton code, feel free to implement any useful functions.

1. enum Node
This class represent a node of type Branch, Leaf, Extension, or Null.

2. struct MerklePatriciaTrie
This class represent a Merkle Patricia Trie. It has two variables: "db" and "root".
Variable "db" is a HashMap. The key of the HashMap is a Node's hash value. The value of the HashMap is the Node. 
Variable "root" is a String, which is the hash value of the root node.

Other requirements:

1. Leaf node and Extension node are differentiated by their prefix, not the enum type. The class "Node" is defined in skeleton code and do not change it!

2. General code quality is required. Code smell would be commented for this project, and might affect your grade in the future projects. 


Hints:

##### 1. Think through all the cases of Get(), Insert(), and Delete() before implementing would save a lot of time. 
##### 2. If you use Rust, implement Clone() for both "Node" and "MerklePatriciaTrie" would be helpful.
##### 3. During the insertion or deletion, if you change the value of an existing node, it's hash_value must be changed. And since the parent node has to update the hash_value so that it could link to that node, the parent's hash_value should be updated as well, all the way up until the root.
##### 4. The grading process includes a lot of test cases, so prepare your own test cases and pass them would help you improve the code. 
##### 5. Given a string, you need to convert it to hex array, and then put that hex array into compact_encode() function to get the encoded hex array.
Here's the logic of converting string to hex array: for each character, calculate the hex format of its ASCII value. For example, given character "a", its ASCII value is 97. The hex format of 97 is 61, then you can put hex array [6, 1] into compact_encode() function. The logic follows the github wixi page, as that page convert "do" to "64 6f", and then put it into compact_encode(). 
##### 6. Normally if we are to implement a tree, we would use "pointers". For example, n1 = Node(), n2 = Node(); then n1.next = n2. But in MerklePatriciaTrie, we use a "HashMap" instead of the "pointer". And we get a node by its hash_value. 
