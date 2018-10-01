# POET Server Specifications
Version for review

## Overview
The POET Server implements the proofs sequential work protocol construction defined in [simple proofs of sequential work](https://eprint.iacr.org/2018/183.pdf). We follow the paper's definitions, construction and are guided by the reference python source code implementation. Please read the paper and analyze the reference python source code. The POET Server is designed to be used by the Spacemesh POET service with is part of the broader Spacemesh protocol but is also useful for other use cases

## Constants
- t:int = 150. A statistical security parameter. (Note: is 150 needed for Fiat-Shamir or 21 suffices?)
- w:int = 256. A statistical security parameter.

Note: The constants are fixed and shared between the Prover and the Verifier. Values shown here are just an example and may be set differently for different POET server instances.

## Input
- n:int - time parameter
- x: {0,1}^w = rndBelow(2^w - 1) - verifier provided input statement

Note: In a real world deployment, n will be constant per POET service instance and known to verfiers using that instance.

## Definitions

- N:int = 2^n-1

- m:int , 0 <= m <= n. Defines how much data should be stored by the prover

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
- We start with Bn - `the complete binary tree of depth n` where all edges go from leaves up the tree to the root, and add edges to the n leaves in the following way
- For each leaf i of the n leaves, we add an edge to the leaf from all the direct siblings of the nodes on the path from the leaf to the root node
- Each node in the DAG is identified by a binary string in the form `0`, `01`, `0010` based on its location in Bn 
- The root node at height 0 is identified by the empty string "" 
- The nodes at height 1 (l0 and l1) are labeled `0` and `1`. The nodes at height 2 are labeled `00`, `01`, `10` and `11`, etc... So for each height h, node's identifier is an h bits binary number that uniquely defines the location of the node in the DAG
- We say node u is a parent of node v if there's a direct edge from u to v in the DAG (based on its construction)
- Each node has a label. The label li of node i (the node with identifier i) is defined as: `li = Hx(i,lp1,...,lpd)` where `(p1,...,pd) = parents(i)`. For example, the root node's label is `lε = Hx("", l0, l1)` as it has 2 only parents l0 and l1 and its identifier is the empty string ""

##### DAG Construction
- Compute the label of each DAG node, and store only the labels of of the dag from the root up to level m
- Computing the labels of the DAG should use up to w * (n + 1) bits of RAM using the following algorithm.

Recursive computation of the labels of DAG(n):

1. Compute the labels of the left subtree (tree with root l0)
2. Keep the label of l0 in memory and discard all other computed labels from memory.
3. Compute the labels of the right subtree (tree with root l1) - using l0.
4. Once l1 is computed, discard all other computed labels from memory and keep l1.
5. Compute the root label le = Hx("", l0, l1)

When a label value is computed by the algorithm, store it in key/value storage if it is in height <= m.

##### DAG Storage
- Please use [LevelDb](https://github.com/syndtr/goleveldb) for storing label values - LevelDB is available as a C++ or a Go lib 
- Labels should be stored keyed by their identifier. e.g. k=i, v= li
- Note hat only up to 1 <= m <= n top layers of the DAG should be stored by POSW(), and the rest should be computed when required on-demand. So storage should size should be O(w * m)
- Use LevelDb caching for fast reads. Cache size should be a verifier param and set based on a deployment runtime available memory settings


