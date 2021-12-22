# Release Notes

## Improvements

- Changed timeout for broadcasting TXs to Umee to 60s to match that of the
  official Gravity Bridge.
- Added a gas limit adjustment flag for Ethereum transactions.

## Bug Fixes

- Use Ethereum gas cost estimation instead of a hardcoded value.
- Claims are split into chunks of 10 to avoid hitting request limits.

## Features

- Add support for fee grant utilization. Operators may populate a `cosmos-fee-granter`
  flag that will allow them to use a fee granter address to pay for relayer fees.
  Note, an existing fee grant must exist between the fee granter and the Peggo
  orchestrator account.
