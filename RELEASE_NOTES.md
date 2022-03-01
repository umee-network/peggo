# Release Notes

This release includes several bug fixes and the ability to check for profit
before requesting a new batch.
It also includes updates to vulnerable dependencies.

## Changelog

### Features

[#216](https://github.com/umee-network/peggo/pull/216) Add profitability check on the batch requester loop.

### Bug Fixes

- [#217](https://github.com/umee-network/peggo/pull/217) Add validation to user input Ethereum addresses.
- [#209](https://github.com/umee-network/peggo/pull/209) Fix the `version` command to display correctly.
- [#205](https://github.com/umee-network/peggo/pull/205) Make sure users are warned when using unencrypted non-local urls in flags.