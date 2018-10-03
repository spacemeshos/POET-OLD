# POET Server Specifications
Version for review

## Overview
The POET Server implements the proofs sequential work protocol construction defined in [simple proofs of sequential work](https://eprint.iacr.org/2018/183.pdf). We follow the paper's definitions, construction and are guided by the reference python source code implementation. Please read the paper and analyze the reference python source code. The POET Server is designed to be used by the Spacemesh POET service with is part of the broader Spacemesh protocol but is also useful for other use cases

## Constants
- t:int = 150. A statistical security parameter. (Note: is 150 needed for Fiat-Shamir or does 21 suffices?). Shared between prover and verifier

- w:int = 256. A statistical security parameter. Shared between prover and verifiers

- n:int - time parameter. Shared between verifier and prover

Note: The constants are fixed and shared between the Prover and the Verifier. Values shown here are just an example and may be set differently for different POET server instances.


## Definitions

- x: {0,1}^w = rndBelow(2^w - 1) - verifier provided input statement (commitment)

- N:int - number of iterations. N := 2^(n+1) - 1. Computed based on n.

- m:int , 0 <= m <= n. Defines how much data should be stored by the prover.

- M : Storage available to the prover, of the form (t + n*t + 1 + 2^{m+1})*w, 0 <= m <= n . For example, with w=256, n=40, t=150 and m=20 prover should use around 70MB of memory and make N/1000 queries for openH.

- Hx : (0, 1)^{<= w(n+1)} => (0, 1)^w . Hx is constructed in the following way: Hx(i) = H(x,i) where H() is a cryptographic hash function. The implementation should use a macro or inline function for H(), and should support a command-line switch that allows it to run with either H()=sha3() or H=sha256().

- φ : (phi) value of the DAG Root label l_epsilon computed by PoSW^Hx(N)

- φP : (phi_P) Result of PoSW^Hx(N) stored by prover. List of labels in the m highest layers of the DAG.

- γ : (gamma) (0,1}^{t*w}. A random challenge sampled by verifier and sent to prover for interactive proofs. Created by concatenation of {gamma_1, ..., gamma_t} where gamma_i = rnd_in_range of (0,1)^w

- τ := openH(x,N,φP,γ) proof computed by prover based on verifier provided challenge γ. A list of t tuples where each tuple is defined as: (l_{gamma_i}, dict{alternate_siblings: l_{the alternate siblings}) for 1 <= i <= t (the security param). So, for each i, the tuple contains the label of the node at index gamma_i, as well as the labels of all siblings of the nodes on the path from the node gamma_i to the root.

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

##### DAG Definitions
- We define n as the depth of the DAG. We set N = 2^(n+1) where n is the time param. e.g. for n=4, N = 31
- We start with Bn - `the complete binary tree of depth n` where all edges go from leaves up the tree to the root, and add edges to the n leaves in the following way
- The DAG has 2^n leaves and 2^n -1 internal nodes
- For each leaf i of the 2^n leaves, we add an edge to the leaf from all the direct siblings of the nodes on the path from the leaf to the root node
- Each node in the DAG is identified by a binary string in the form `0`, `01`, `0010` based on its location in Bn. This is the node id.
- The root node at height 0 is identified by the empty string ""
- The nodes at height 1 (l0 and l1) are labeled `0` and `1`. The nodes at height 2 are labeled `00`, `01`, `10` and `11`, etc... So for each height h, node's id is an h bits binary number that uniquely defines the location of the node in the DAG
- We say node u is a parent of node v if there's a direct edge from u to v in the DAG (based on its construction)
- Each node has a label. The label li of node i (the node with id i) is defined as:

```
li = Hx(i,lp1,...,lpd)` where `(p1,...,pd) = parents(i)
```

For example, the root node's label is `lε = Hx("", l0, l1)` as it has 2 only parents l0 and l1 and its id is the empty string "".

##### Computing node parents ids
Given a node i in a dag(n), we need a way determine its set of parent nodes. For example, we use the set to compute its label. This can be implemented without having to store all DAG edges in storage.

Note that with this binary string labeling scheme we get the following properties:

1. The id of left sibling of a node in the dag is node i label with the last bit flipped from 1 to 0. e.g. the left sibling of node with id `1001` is `1000`
2. The id of a direct parent in Bn of a node i equals to i with the last bit removed. e.g. the parent of node with id `1011` is `101`

- Using these properties, the parents ids can be computed based only on the DAG definition and the node's identifier by the following algorithm:

`If id has n bits (node is a leaf in dag(n)) then add the ids of all left siblings of the nodes on the path from the node to the root, else add to the set the 2 nodes below it (left and right nodes) as defined by the binary tree Bn.`

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

##### DAG Construction
- Compute the label of each DAG node, and store only the labels of of the dag from the root up to level m
- Computing the labels of the DAG should use up to w * (n + 1) bits of RAM using the following algorithm.

Recursive computation of the labels of DAG(n):

1. Compute the labels of the left subtree (tree with root l0)
2. Keep the label of l0 in memory and discard all other computed labels from memory
3. Compute the labels of the right subtree (tree with root l1) - using l0
4. Once l1 is computed, discard all other computed labels from memory and keep l1
5. Compute the root label le = Hx("", l0, l1)

- When a label value is computed by the algorithm, store it in a key/value storage if the label's height <= m.
- Note that this works because only l0 is needed for computing labels in the tree rooted in l1. All of the additional edges to nodes in the tree rooted at l1 start at l0.
- Note that the reference Python code does not construct the DAG in this manner and keeps the whole DAG in memory. Please use the Python code as an example for simpler constructions such as binary strings, open and verify.

##### DAG Storage
- Please use [LevelDb](https://github.com/syndtr/goleveldb) for storing label values - LevelDB is available as a C++ or a Go lib
- Labels should be stored keyed by their id. e.g. k=i, v= li
- Note hat only up to 1 <= m <= n top layers of the DAG should be stored by POSW(), and the rest should be computed when required on-demand. So storage should size should be O(w * m)
- Use LevelDb caching for fast reads. Cache size should be a verifier param and set based on a deployment runtime available memory settings

##### APIs

// See notes about data types below (commitment, proof, challenge)

```
Verifier {
    // Set new commitment and provide callback for POET server result POSW(n)
    // Verifer should start a new prover with the provided commitment and n
    // The callback includes the NIP proof or an error.
    SetCommitment(commitment: bytes, n: int, callback: (proof: Proof, error));
    
    // Verify a proof
    Verify(proof);
    
    // Verify a random challenge
    VerifyRandomChallenge() returns (result:bool, error: Error);
}

Prover {
    // start POSW(n) and return NIP or error in callback after POSW(n) is complete
    Start(commitment: bytes, n: int, callback: (result: Proof, error: Error);
    
    // returns a proof based on challenge
    GetProof(challenge: Challenge);
}

TestNip() {
    const n = 40;
    const c = randomBytes(32)
    v = new Verifier();
    v.SetCommitment(c, n callback);
    
    callback(result: Proof, error: Error) {
        assertNoError(error);    
        res = v.Verify(result);
        assert(res);
    }
}

TestBasicRandomChallenge() {
    const n = 40;
    const c = crypto.randomBytes(32)
    v = new Verifier();
    
    v.SetCommitment(c, n, callback);
    callback(result: proof, error: Error) {
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
    callback(result: Proof, error: Error) {
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
arbitray length bytes. Verifier implementation should just use H(commitment) to create a commitment that is in range (0, 1)^w . So the actualy commitment can be sha256(commitment) when w=256.

#### Challenge 
A challenge is a list of t random binary strings in {0,1}^n. Each binary string is an identifier of a leaf node in the DAG.
Note that the binary string should always be n bytes long, including trailing `0`s if any, e.g. `0010`.

#### Proof
A proof needs includes the following data:
1. φ - the label of the root node.
2. For each identifier i in a challange (0 <= i < t), an ordered list of labels which includes: 
   2.1 li - The label of the node i
   2.2 An ordered list of the labels of the sibling node of each node on the path to the parent node.

So, for example for Dag(4) and for a challenge identifier `0101` - The labels that should be included in the list are: 0101, 0100, 011, 00 and 1. This is basically an opening of a merkle tree commitment.

The complete proof data can be encoded in a tuple where the first value is φ and the second value is a dictionary with an entry for each of the t challenge identifiers using the following syntax:

{ φ, { identifier_0 : { label(0), list_of_siblings_on_path_to_root_from_0}, .... { identifier_t : { label(t), list_of_siblings_on_path_to_root_from_t} }






