# Changelog

## [0.0.2] - 2024-01-02

### Bug Fixes

- Reformat go.mod
- Change argument number condition
- Reformat encrypt decrypt flag reading

### Features

- Store flags in viper instead of global variables
- Bind root command to function
- Add secret-key and secret-file flags
- Make config file type dynamic (json, yaml, toml)
- Bind flag to env
- Implement age encryption
- Update go.mod and go.sum
- Add github actions workflow for changelog and binary
- Add util function for getting file size and mode
- Support decryption
- Handle encrypt and decrypt flag
- Add config struct
- Combine ageEncrypt with encrypt and ageDecrypt with decrypt
- Implement storage upload and download
- Implement gcs upload and download
- Implement s3 upload and download
- Add more utility for file handling and naming
- Update go.mod and go.sum

### Refactor

- Improve age encryption logic
## [0.0.1] - 2023-12-09

### Bug Fixes

- Change config declaration syntax
- Remove .yaml on config error log

### Features

- Add base cli
- Fork hashicorp shamir lib
- Add minimum number of args
- Change encrypt and decrypt command to options
- Add changelog
