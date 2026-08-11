# Repository instructions

- The module path is `github.com/btnguyen2k/cf-base`.
- Maintain compatibility with Go 1.18 and later. Do not use language features or dependencies that require a newer Go version.
- Place tests in `module_test/`, which is maintained as a separate Go module.
- Tests that need access to unexported identifiers may instead be placed beside the production code and use the same package.
