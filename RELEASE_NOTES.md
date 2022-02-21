# Release Notes

This release introduces the deprecation of the `--relay-valsets` flag in favor
of `--valset-relay-mode` which allows a finer control over how valsets will be
relayed.
The ERC20 to Coingecko IDs mapping was also updated to accomodate the new ERC20s
deployed this week.

## Changelog

### Features

- [#189] Add the flag `--valset-relay-mode` which allows a finer control over
  how valsets will be relayed.

### Improvements

- [#201] Add ERC20 mappings for Umee's new tokens.

### Deprecated

- [#189] Deprecate the `--relay-valsets` flag.