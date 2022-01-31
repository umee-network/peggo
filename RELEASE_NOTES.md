# Release Notes

This release introduces bug fixes and improvements.

## Changelog

### Improvements

- [#132] Add `cosmos-msgs-per-tx` flag to set how many messages (Ethereum claims)
  will be sent in each Cosmos transaction.
- [#134] Improve valset relaying by changing how we search for the last valid
  valset update.

### Bug Fixes

- [#134] Fix logs, CLI help and a panic when a non-function call transaction was
 received during the TX pending check.