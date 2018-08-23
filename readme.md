## POET SERVICE
Draft

### Overview
POET stands for proofs of elapsed time. It is aka proofs of sequential work.
The Spacemesh POET service is a public Internet service that provides verifiable proofs of elapsed time. It is designed to work together with Spacemesh proofs of space (POST components) commitments to provide a `NIPST` - a non-interactive proof of space time.

THIS IS A PRE-REVIEW PUBLIC DRAFT - comments and suggestions are welcomed.

## Requirements

- A POET service must be configured for a specific time duration number of cycles which roughly corresponds to a time duration based on the deployed POET service CPU single core frequency. Or more accurately put, to the performance of the underlying hash function used by the service on the deployed hardware.
- The POET must provide a service that allows anyone to create a NIPST for a commitment.
- The POET service should be configured with a public time beacon. We plan to use the timestamp of irreversible Spacemesh blockmesh layers as the time beacon. The time beacon guarantees that a POET proof is not older than the published layer timestamp.
- The POET service should be provided as an `https-json` for any client, and as an `gRpc` endpoint for gRPC clients.
- The POET service should work in consecutive rounds. Clients submitted statements which are received before a round starts must be part of the initial round commitment. In other words, a POET round initial statement must be a proof of all received statements before the round start.  

## Solution Design
- Our design follows the theoretical work of Tal Moran's [Publicly verifiable proofs of sequential work](https://eprint.iacr.org/2011/553.pdf) using the data structure optimizations published in [Simple Proofs of Sequential Work](https://eprint.iacr.org/2018/183.pdf)

### POET Service Config
- A POET service is configured with one ore more `Spacemesh API gateways`. Each gateway provides the `Spacemesh API` to the Spacemesh mainent. The service uses the Spacemesh API to obtain irreversible layers meta-data such as hash and timestamp.
- Each service should have a crypto key pair used to sign statements created by the service and for anyone to verify statements signed by the service.
- Round 0 of a Spacemesh service will not use a Spacemesh layer id as it is designed to provide the initial proof required to select validators for a Spacemesh network
- A round may fail due to a runtime server error. The service should report failed rounds.

### POET SERVICE API

- `GetServiceInfo`
    - Public Key
    - Current round id (counter)
    - Last completed round id (counter)


- `GetRoundInfo(roundId)`
    Response:
    - RoundId (counter)
    - Service public key
    - Round start timestamp
    - Round status: running / failed / complete
    - Spacmesh blockmesh {id, hash, timestamp} - when applicable (roundId > 0)
    - Round complete timestamp (n/a for running rounds)
    - Round Manifest:
        - Ordered list of client submitted statements
        - Signature on the list
        - Ordered list hash - this is used as the x in PoSWHx()
    - proof - a non-interactive proof for a completed round.
    - Signature: service public id signature on the proof


- `SubmitStatement(data)`
    - data: binary data
    Response:
        - ok / error


### Theoretical background
- https://eprint.iacr.org/2011/553.pdf
- https://eprint.iacr.org/2018/183.pdf

### Related work
- https://github.com/wfus/proof-of-sequential-work
