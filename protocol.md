# POET Server Specifications
Draft

## Overview
The POET Server implements the proofs sequential work protocol construction defined in [simple proofs of sequential work](https://eprint.iacr.org/2018/183.pdf). We follow the paper's definitions, construction and are guided by the reference python source code implementation. Please read the paper and analyze the reference python source code. The POET Server is designed to be used by the Spacemesh POET service with is part of the broader Spacemesh protocol but is also useful for other use cases.

Section numbers in "Simple proofs of sequential work" are referenced by this spec.


## Constants (See section 1.2)
- t:int = 150. A statistical security parameter. (Note: is 150 needed for Fiat-Shamir or does 21 suffices?). Shared between prover and verifier

- w:int = 256. A statistical security parameter. Shared between prover and verifiers

- n:int - time parameter. Shared between verifier and prover

Note: The constants are fixed and shared between the Prover and the Verifier. Values shown here are just an example and may be set differently for different POET server instances.


## Definitions (See section 4, 5.1 and 5.2)

- x: {0,1}^w = rndBelow(2^w - 1) - verifier provided input statement (commitment)

- N:int - number of iterations. N := 2^(n+1) - 1. Computed based on n.

- m:int , 0 <= m <= n. Defines how much data should be stored by the prover.

- M : Storage available to the prover, of the form (t + n*t + 1 + 2^{m+1})*w, 0 <= m <= n . For example, with w=256, n=40, t=150 and m=20 prover should use around 70MB of memory and make N/1000 queries for openH.

- Hx : (0, 1)^{<= w(n+1)} => (0, 1)^w . Hx is constructed in the following way: Hx(i) = H(x,i) where H() is a cryptographic hash function. The implementation should use a macro or inline function for H(), and should support a command-line switch that allows it to run with either H()=sha3() or H=sha256().

- φ : (phi) value of the DAG Root label l_epsilon computed by PoSW^Hx(N)

- φP : (phi_P) Result of PoSW^Hx(N) stored by prover. List of labels in the m highest layers of the DAG.

- γ : (gamma) (0,1}^{t*w}. A random challenge sampled by verifier and sent to prover for interactive proofs. Created by concatenation of {gamma_1, ..., gamma_t} where gamma_i = rnd_in_range of (0,1)^w

- τ := openH(x,N,φP,γ) proof computed by prover based on verifier provided challenge γ. A list of t tuples, where each tuple is defined as follows for 0 <= i < t: {label_i, list_of_siblings_to_root_from_i }. So, for each i, the tuple contains the label of the node width identifier i, as well as the labels of all siblings of the nodes on the path from the node i to the root.

- NIP (Non-interactive proof): a proof τ created by computing openH for the challenge γ := (Hx(φ,1),...Hx(φ,t)). e.g. without receiving γ from the verifier. Verifier asks for the NIP and verifies it like any other openH using verifyH.

- verifyH(x,N,φ,γ,τ) ∈ {accept, reject} - function computed by verifier based on prover provided proof τ.

### Actors
- Prover: The service providing proofs for verifiers
- Verifier: A client using the prover by providing the input statement x, and verifying the prover provided proofs (by issuing random challenges or by verifying a non-interactive verifier provided proof for {PoSW^Hx(N), x}

## Base Protocol Test Use Cases

### User case 1 - Basic random challenges verification test
1. Verifier generates random commitment x and sends it to the prover
2. Prover computes PoSWH(x,N) by making N sequential queries to H and stores φP
3. Prover sends φ to verifier
4. Verifier creates a random challenge γ and sends it to the prover
5. Prover returns proof τ to the verifier
6. Verifier computes verifyH() for τ and outputs accept

### Use case 2 - NIP Verification test
1. Verifier generates random commitment x and sends it to the prover
2. Prover computes PoSWH(x,N) by making N sequential queries to H and stores φP
3. Prover creates NIP and sends it to the verifier
4. Verifier computes verifyH() for NIP and outputs accept

### User case 3 - Memory requirements verification
- Use case 1 with a test that prover doesn't use more than M memory

### User case 4 - Bad proof detection
- Modify use case 1 where a random bit is changed in τ proof returned to the verifier by the prover
- Verify that verifyH() outputs reject for the verifier

### User case 5 - Bad proof detection
- Modify use case 2 where a random bit is changed in the NIP returned to the verifier by the prover
- Verify that verifyH() outputs reject for the verifier

### Theoretical background and context
- [1] https://eprint.iacr.org/2011/553.pdf
- [2] https://eprint.iacr.org/2018/183.pdf
- [3] https://spacemesh.io/whitepaper1/
- [4] https://pdfs.semanticscholar.org/b904/6d002da153a6fe9b06d469da4efffdfcb9c6.pdf

### Related work
- [5] https://github.com/avive/proof-of-sequential-work
- [6] https://github.com/avive/slow-time-functions

### Implementation Guidelines

#### DAG
The core data structure used by the verifier.

##### DAG Definitions (See section 4)
- We define n as the depth of the DAG. We set N = 2^(n+1) where n is the time param. e.g. for n=4, N = 31
- We start with Bn - `the complete binary tree of depth n` where all edges go from leaves up the tree to the root. Bn has 2^n leaves and 2^n -1 internal nodes. We add edges to the n leaves in the following way:

    For each leaf i of the 2^n leaves, we add an edge to the leaf from all the direct siblings of the nodes on the path from the leaf to the root node.

    In other words, for every leaf u, we add an edge to u from node v_{b-1}, iff v_{b} is an ancestor of u and nodes v_{b-1}, v_{b} are direct siblings.

- Each node in the DAG is identified by a binary string in the form `0`, `01`, `0010` based on its location in Bn. This is the node id
- The root node at height 0 is identified by the empty string ""
- The nodes at height 1 (l0 and l1) are labeled `0` and `1`. The nodes at height 2 are labeled `00`, `01`, `10` and `11`, etc... So for each height h, node's id is an h bits binary number that uniquely defines the location of the node in the DAG
- We say node u is a parent of node v if there's a direct edge from u to v in the DAG (based on its construction)
- Each node has a label. The label li of node i (the node with id i) is defined as:

```
li = Hx(i,lp1,...,lpd)` where `(p1,...,pd) = parents(i)
```

For example, the root node's label is `lε = Hx(bytes(""), l0, l1)` as it has 2 only parents l0 and l1 and its id is the empty string "".

##### Implementation Note: packing values for hashing

- To pack an identifier e.g. "00111" value for the input of Hx(), encoded it as a byte array as utf-8 bytes array. For example, in Go use: []byte("00011")
- Labels are arbitrary 32 bytes of binary data so they don't need any encoding.
- As an example, to compute the input for Hx("001", label1, label2), encode the binary string to a utf-8 encoded bytes array and append to it the labels byte arrays.



##### Computing node parents ids
Given a node i in a DAG(n), we need a way determine its set of parent nodes. For example, we use the set to compute its label. This can be implemented without having to store all DAG edges in storage.

Note that with this binary string labeling scheme we get the following properties:

1. The id of left sibling of a node in the DAG is node i label with the last (LSB) bit flipped from 1 to 0. e.g. the left sibling of node with id `1001` is `1000`
2. The id of a direct parent in Bn of a node i equals to i with the last bit removed. e.g. the parent of node with id `1011` is `101`

- Using these properties, the parents ids can be computed based only on the DAG definition and the node's identifier by the following algorithm:

`If id has n bits (node is a leaf in DAG(n)) then add the ids of all left siblings of the nodes on the path from the node to the root, else add to the set the 2 nodes below it (left and right nodes) as defined by the binary tree Bn.`

- So for example, for n=4, for the node l1 with id `0`, the parents are the nodes with ids `00` and `01` and the ids of the parents of leaf node `0011` are `0010` and `000`. The ids of the parents of node `1101` are `1100`, `10` and `0`.

- The following Python function demonstrates how to implement this algorithm. It returns a sorted set of parent ids for input which consists of node id (binary string) and the value of n (int):

``` python
def get_parents(binary_str, n=DEFAULT_n):
    parents = []
    bit_list = binary_str.get_bit_list()
    length = binary_str.length
    if length == n:
        for i in range(1, length + 1):
            if bit_list[i - 1] == 1:
                parents.append(BinaryString(i, bits_to_int(bit_list[:i - 1] + [0])))
    else:
        parents.append(BinaryString(length + 1, bits_to_int(bit_list + [0])))
        parents.append(BinaryString(length + 1, bits_to_int(bit_list + [1])))
    return sorted(parents)
```

- Note that leaf nodes parents are not all the siblings on the path to the root from the leaf. The parents are only the left siblings on that path. e.g. siblings with ids that end with a `0`.

##### DAG Construction (See section 4, Lemma 3)
- Compute the label of each DAG node, and store only the labels of of the DAG from the root up to level m
- Computing the labels of the DAG should use up to w * (n + 1) bits of RAM
- The following is a possible algorithm that satisfies these requirements. However, any implementation that satisfies them (with equivalent or better computational complexity) is also acceptable.

Recursive computation of the labels of DAG(n):

1. Compute the labels of the left subtree (tree with root l0)
2. Keep the label of l0 in memory and discard all other computed labels from memory
3. Compute the labels of the right subtree (tree with root l1) - using l0
4. Once l1 is computed, discard all other computed labels from memory and keep l1
5. Compute the root label le = Hx("", l0, l1)


- When a label value is computed by the algorithm, store it in persistent storage if the label's height <= m.
- Note that this works because only l0 is needed for computing labels in the tree rooted in l1. All of the additional edges to nodes in the tree rooted at l1 start at l0.
- Note that the reference Python code does not construct the DAG in this manner and keeps the whole DAG in memory. Please use the Python code as an example for simpler constructions such as binary strings, open and verify.

##### DAG Storage
Please use a binary data file to store labels and not a key/value db. Labels can be stored in the order in which they are computed.

Given a node id, the index of the label value in the data file can be computed by:

     idx = sum of sizes of the subtrees under the left-siblings on path to root + node's own subtree.

The size of a subtree under a node is simply `2^{height+1}-1` * the label size. e.g. 32 bytes.

##### APIs


// See notes about data types below (Commitment, Proof and Challenge)

```
Base { // shared functionality and constants between verifier and Prover

    // Returns the NIP's challenge based on the provided x and φ - phi  
    // Implementation: c = (Hx(φ,1),...Hx(φ,t))
    CreteNipChallenge(x: bytes, phi: bytes) returns c:Challenge;
}

Verifier extends Base {
    // Set new commitment and provide callback for POET server result POSW(n)
    // Verifier should start a new prover with the provided commitment and n
    // The callback includes the NIP proof or an error.
    SetCommitment(commitment: bytes, n: int, callback: (proof: Proof, error));

    // Verify a proof for a challenge
    Verify(challenge: Challenge, proof: Proof);

    // Verify a random challenge
    VerifyRandomChallenge() returns (result:bool, error: Error);
}

Prover extends Base {
    // start POSW(n) and return NIP or error in callback after POSW(n) is complete
    Start(commitment: bytes, n: int, callback: (result: Proof, error: Error);

    // returns a proof based on challenge
    GetProof(challenge: Challenge) returns Proof;
}

// This method implements verifyH(x,N,φ,γ,τ) using provided input arguments and
// verifier constants
Verifer.Verify(challenge: Challenge, proof: Proof) {

    // Verify that the sets of identifiers in challenge and proof are identical.

    phi = proof.phi;
    for (i=0; i < t; i++) {

        // Note that verifier knows the identifier from the challenge it issued and can't trust the prover
        // to return the correct id. So node_id is taken from the challenge and not from the proof:
        node_id =  challenge.identifer(i);

        node_label = proof.label(i);

        // Note that labels that were already included in sibling lists in the proofs will be omitted from the sibling list. The verifier should to use the labels of node ids it already knows about. e.g. It builds an in-memory dictionary [node_id : node_lablel] and update it as a proof data is read. If the dictionary includes an entry for a node_id then the value should be read from the dictionary and not from the sibling list.

        siblings = proof.siblings(i);

        // Check the Merkle-like commitment
        // Compute the labels of the nodes on the path from node_id to the root node phi
        // using the siblings labels and node_label
        // return false if the computed root node does not equal to phi.
    }

    return true;
}

Verifier.VerifyRandomChallenge() returns bool {

    // Create a valid challenge with t random leaf identifiers
    challenge = new Challenge();

    // Get a proof from the verifier for the challenge
    proof = verifier.GetProof(challenge);

    // verify the proof
    return verify(challenge, proof);
}

TestNip() {
    const n = 40;
    const c = randomBytes(32)
    v = new Verifier();
    v.SetCommitment(c, n callback);

    callback(proof: Proof, error: Error) {
        assertNoError(error);    
        c: Challenge = CreteNipChallenge();
        res = v.Verify(c, proof);
        assert(res);
    }
}

TestBasicRandomChallenge() {
    const n = 40;
    const c = crypto.randomBytes(32)
    v = new Verifier();

    v.SetCommitment(c, n, callback);
    callback(proof: Proof, error: Error) {
        assertNoError(error)
        res = v.VerifyRandomChallenge();
        assert(res);
    }
}

TestRndChallenges() {
    const n = 40;
    const c = randomBytes(32)
    v = new Verifier();
    v.SetCommitment(c, n, callback);
    callback(proof: Proof, error: Error) {
        assertNoError(error);
        for (i = 1 to 1000) {
            res = v.VerifyRandomChallenge();
            assert(res);
        }
    }
}


```

### Data Types

#### Commitment
arbitrary length bytes. Verifier implementation should just use H(commitment) to create a commitment that is in range (0, 1)^w . So the actually commitment can be sha256(commitment) when w=256.

#### Challenge
A challenge is a list of t random binary strings in {0,1}^n. Each binary string is an identifier of a leaf node in the DAG.
Note that the binary string should always be n bytes long, including trailing `0`s if any, e.g. `0010`.

#### Proof (See section 5.2)
A proof needs includes the following data:
1. φ - the label of the root node.
2. For each identifier i in a challenge (0 <= i < t), an ordered list of labels which includes:
   2.1 li - The label of the node i
   2.2 An ordered list of the labels of the sibling node of each node on the path to the parent node, omitting siblings that were already included in previous siblings list int he proof

So, for example for DAG(4) and for a challenge identifier `0101` - The labels that should be included in the list are: 0101, 0100, 011, 00 and 1. This is basically an opening of a Merkle tree commitment.

The complete proof data can be encoded in a tuple where the first value is φ and the second value is a list with t entries. Each of the t entries is a list starting with the node with identifier_t label, and a node for each sibling on the path to the root from node identifier_t:

{ φ, {list_of_siblings_on_path_to_root_from_0}, .... {list_of_siblings_on_path_to_root_from_t} }

Note that we don't need to include identifier_t in the proof as the identifiers needs to be computed by the verifier.

Also note that the proof should omit from the siblings list labels that were already included previously once in the proof. The verifier should create a dictionary of label values keyed by their node id, populate it from the siblings list it receives, and use it to get label values omitted from siblings lists - as the verifier knows the ids of all siblings on the path to the root from a given node.


### About NIPs

a NIP is a proof created by computing openH for the challenge γ := (Hx(φ,1),...Hx(φ,t)). e.g. without receiving γ from the verifier. Verifier asks for the NIP and verifies it like any other openH using verifyH. Note that the prover generates a NIP using only Hx(), t (shared security param) and φ (generated by PoSW(n)). To verify a NIP, a verifier generates the same challenge γ and verifies the proof using this challenge.

Hx(φ,i) output is `w bits` but each challenge identifier should be `n bits` long. To create an identifier from Hx(φ,i) we take `the leftmost t bits` - starting from the most significant bit. So, for example, for n == 3 and for Hx(φ,1) = 010011001101..., the identifier will be `010`.

### GetProof(challenge: Challenge) notes
The verifier only stores up to m layers of the DAG (from the root at height 0, up to height m) and the labels of the n leaves.
Generating a proof involves computing the labels of the siblings on the path from a leaf to the DAG root where some of these siblings are in DAG layers with height > m. These labels are not stored by the prover and need to be computed when a proof is generated. The following algorithm describes how to compute these siblings:

1. For each node_id included in the input challenge, compute the node id of the node n. The node on the path from node_id to the root at DAG level m

2. Construct the DAG rooted at node n. When the label of a sibling on the path from node_id to the root is computed as part of the DAG construction, add it to the list of sibling labels on the path from node_id to the root

### Computing Hx(y,z) - Argument Packing
When we hash 2 (or more) arguments to hash, for example Hx(φ,1).
We need to agree on a canonical way to pack params across implementations and tests into 1 input bytes array.
We define argument packing in the following way:

- Each argument is serialized to a []byte.
- String serialization: utf-8 encoding to a byte array
- Int serialization (uint8, ... uint64): BigEndian encoding to a byte array
- The arguments are concatenated by appending the bytes of each argument to the bytes of the previous argument. For example, `Hx(φ,i) = Hx([]byte(φ) ... []byte(i))`
