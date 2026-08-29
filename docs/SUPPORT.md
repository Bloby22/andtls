# Support

## Getting Help

- **Documentation**: See [README.md](../README.md) and [docs/](./)
- **Issues**: [GitHub Issues](https://github.com/Bloby22/andtls/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Bloby22/andtls/discussions)

## Troubleshooting

### ADB not found

Ensure `adb` is installed and available in your `PATH`. Install Android SDK Platform-Tools from https://developer.android.com/tools/releases/platform-tools.

### No devices detected

1. Enable **USB Debugging** on your Android device (Settings → Developer Options)
2. Check the USB cable connection
3. Accept the RSA key prompt on the device
4. Run `adb devices` to verify

### Wireless ADB not working

- Ensure both the computer and device are on the same WiFi network
- The device must have USB debugging enabled initially
- TCP mode uses port **5555** by default

### Screenshots not saving

Check that the `./screenshots/` directory is writable and exists.

## Community

Please read our [Code of Conduct](./CODE_OF_CONDUCT.md) before participating.
