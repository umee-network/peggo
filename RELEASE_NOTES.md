# Release Notes

This release fixes one of the main issues from the Trail of Bits audit, the need
for multiple price oracles instead of one. Now we can have multiple, thanks to
the implementation of the price feeder as a module.

## Changelog

### Features

[#231](https://github.com/umee-network/peggo/pull/231) Add multiple providers for token prices.