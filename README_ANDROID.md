# Android

KeepassUI can be distributed as an Android App as well.

## Android Distribution

### Install Android tools

To install the required android tools to build the apk and bundle file:

```make install-all-android-tools```

### Generate APK and AAB files

To generate an Android APK:

```make package-android```

To create a KeyStore file required for the next step:

```keytool -genkey -v -keystore gplay.keystore -alias alias -keyalg RSA -keysize 2048 -validity 9999```

To generate an Android Bundle to upload to the Google Play Store:

```make release-android``` (requires the presence of the keystore file `~/dev/gplay.keystore`.)

