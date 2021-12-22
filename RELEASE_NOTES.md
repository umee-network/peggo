# Release Notes

## Bug Fixes

- Use Ethereum gas cost estimation instead of a hardcoded value.

## Features

- Add support for fee grant utilization. Operators may populate a `cosmos-fee-granter`
  flag that will allow them to use a fee granter address to pay for relayer fees.
  Note, an existing fee grant must exist between the fee granter and the Peggo
  orchestrator account.
