# Security Policy

## Supported Versions

| Version | Supported          |
|---------|--------------------|
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability within andtls, please send an email to **blobycz@gmail.com**. All security vulnerabilities will be promptly addressed.

**Please do not report security vulnerabilities through public GitHub issues.**

### What to include

- Description of the vulnerability
- Steps to reproduce the issue
- Affected versions
- Any potential impact assessment

### Response timeline

- **Acknowledgment**: within 48 hours
- **Initial assessment**: within 1 week
- **Fix or mitigation**: depends on severity

## Security Considerations

andtls interacts with Android devices via ADB. Users should be aware:

- ADB connections should only be made over trusted networks
- USB debugging should be disabled when not in use
- Wireless ADB (TCP mode) exposes the device on the network — use with caution
