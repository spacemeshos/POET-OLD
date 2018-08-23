## POET SERVICE
Draft

### Overview
POET stands for proofs of elapsed time. It is aka proofs of sequential work.
The Spacemesh POET service is a public Internet service that provides verifiable proofs of elapsed time. It is designed to work together with Spacemesh proofs of space (POST components) commitments to provide a `NIPST` - a non-interactive proof of space time.

## Requirements

### Service Configuration
- A POET service must be configured for a specific time duration number of cycles which roughly corresponds to a time duration based on the deployed POET service CPU single core frequency. Or more accurately put, to the performance of the underlying hash function used by the service on the deployed hardware.
- The POET must provide a service that allows anyone to create a NIPST for a commitment.
- A POET service should be configured with a public time beacon. We plan to use the timestamp of irreversible Spacemesh blockmesh layers as the time beacon. The time beacon guarantees that a POET proof is not older than the published layer timestamp.

## Design
### POET API

- `GetServiceInfo`

- `GetRoundInfo(roundId)`

- `SubmitStatement(statementBinaryData)`
