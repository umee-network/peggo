# Release Notes

This release introduces the deprecation of an unsafe flag and some helpers for
the Coingecko price feed.

## Changelog

### Improvements

- [#172] Add fallback token addresses (to aid price lookup)
- [#185] Add fallback token addresses (to aid price lookup) for Umee

### Deprecated

- [#174] Deprecate `--eth-pk` in favor of an env var (`$PEGGO_ETH_PK`)