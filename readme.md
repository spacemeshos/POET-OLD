## POET SERVICE
Draft

### Overview
- POET stands for Proofs Of Elapsed Time. It is a.k.a. proofs of sequential work. The Spacemesh POET service is a public Internet service that provides verifiable proofs of elapsed time
- It is designed to work together with Spacemesh Proofs Of Space Time (POST) commitments to provide NIPSTs - non-interactive proofs of space time

THIS IS A PRE-REVIEW PUBLIC DRAFT

#### Open Design Process
- We are designing the POET service fully in the open before moving on to the open source code implementation
- You are welcome to join the design phase. Collaborators and contributors are welcomed
- To get started join our [POET Gitter channel](https://gitter.im/spacemesh-os/POET) and introduce yourself

## Service Requirements
- A POET service participates in the the creation of NIPSTs (Non Interactive Proofs of Space Time) as specified in the [Spacemesh White Paper 1](https://spacemesh.io/whitepaper1/)
- A POET service must be configured for a number of iterations which roughly corresponds to a time duration based on the deployed POET service CPU single core frequency. Or more accurately put, to the performance of the underlying hash function used by the service on the deployed hardware
- The POET must provide a service that allows anyone to create a NIPST for a commitment
- The POET service should be configured with a public time beacon. We plan to use the timestamp of irreversible Spacemesh blockmesh layers as the time beacon. The time beacon guarantees that a POET proof is not older than the published layer timestamp
- The POET service should be provided as an `https-json` for any client, and as an `gRpc` endpoint for gRpc clients
- The POET service should work in consecutive rounds. Clients submitted statements which are received before a round starts must be part of the initial round commitment. In other words, a POET round initial statement must be a proof of all received statements before the round start
- A round may fail due to a runtime server error. The service should report failed rounds

## Solution Design
- Our design follows the theoretical work of Tal Moran's [Publicly verifiable proofs of sequential work](https://eprint.iacr.org/2011/553.pdf) using the data structure optimizations published in [Simple Proofs of Sequential Work](https://eprint.iacr.org/2018/183.pdf)
- The design will be reviewed by the Spacemesh research team before implementation begins
- We plan to initially use sha256 as the bash H hash function 
- We plan to implement the service in go-lang to achieve close to native pref and native cross-platform packaging
- The statement X used in each round to generate Hx(s):=sha-256(X||s) is defined as:

    `x := {Service signature on the hash of the client submitted statements sorted list || RoundId || Spacemesh blockmesh layer hash}` where || signifies binary concatenation.
    
    - Spacemesh blockmesh layer hash: the hash of the irreversible layer used for the round. This proves that the POET proof was started after the layer timestamp.
    - Client submitted statements sorted list - the statements submitted in time for participation in this round signed by the service to prove that he used them to generate the poet proof
    - RoundId - current round id

### POET Service Config
- A POET service is configured with one ore more `Spacemesh API gateways`. Each gateway provides the `Spacemesh API` to the Spacemesh mainent. The service uses the Spacemesh API to obtain irreversible layers meta-data such as hash and timestamp
- Each service should have a crypto key pair used to sign statements created by the service and for anyone to verify statements signed by the service
- Round 0 of a POET service will not use a Spacemesh layer id as it is designed to provide the initial proof required to select validators for a Spacemesh network. For round 0 a POET service will be confiured with another form of time beacon. e.g. a recent hash of another public blockchain block to ensure that the proof was created after the block timestamp wall clock

### POET Service Api

#### API Methods

- `GetServiceInfo`
    - Service public key
    - Currently running round id (counter)
    - Last completed round id (counter)


- `GetRoundInfo(roundId)`
    Response:
    - RoundId (counter)
    - Round status: running / failed / complete
    - Service public key
    - Round start timestamp
    - Spacmesh blockmesh layer data {id, hash, timestamp} - when applicable (roundId > 0)
    - Round complete timestamp (n/a for running rounds)
    - Round Manifest:
        - Ordered list of client submitted binary statements
        - Poet signature of the list
        - Ordered list hash - this is used as the x in PoSWHx()
    - Proof: a non-interactive proof for a completed round (Using the Fiat-Shamir heuristic)
    - Signature: service public id signature of the proof
    - X: X binary date used for Hx() - verifiers need it to verify the proof


- `SubmitStatement(data)`
    - data: arbitrary binary data (with a hard-coded limit on the number of bytes)
    Response:
        - status: ok / error
        - expectedRoundId: Id of round that a proof will be provided for this statement. (This will simply be the id of the currently running round plus 1).

### Implementation Considerations
- The service should be implemented as an https-json service with json as both the request params and response data format
- All binary data should be `base64` encoded in json payloads
- We plan using a modern, optimized implementation of sha-256 for modern Intel CPUs. See: https://github.com/avive/slow-time-functions


### Theoretical background and context
- [1] https://eprint.iacr.org/2011/553.pdf
- [2] https://eprint.iacr.org/2018/183.pdf
- [3] https://spacemesh.io/whitepaper1/

### Related work
- https://github.com/avive/proof-of-sequential-work
- https://github.com/avive/slow-time-functions
