![License](https://img.shields.io/github/license/npalumbo/keepassui) ![Build Status](https://github.com/npalumbo/keepassui/actions/workflows/run_tests.yaml/badge.svg)  ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/npalumbo/keepassui/main) 

# KeepassUI
A password manager port of [Keepass](https://keepass.info/). Built using the [gokeepasslib](https://github.com/tobischo/gokeepasslib) and [fyne](https://github.com/fyne-io/fyne) libraries.

## Motivation
I noticed that some of the commercial password managers don't support desktop and mobile versions at the same time on their free versions. We can solve that problem using [Keepass](https://keepass.info/). 
Additionally, Keepass has [many ports](https://keepass.info/download.html), many of them to specific target platforms like Mac or Android. I thought it was interesting to explore writing a Keepass UI that can be released on many platforms from the same codebase!.

## Goal
The desired state is for this to evolve into a desktop and mobile app using the same codebase. This should be achieved thanks to the Golang [fyne](https://github.com/fyne-io/fyne) GUI toolkit.

## Contributing
This is in early stages. All help is welcome. If you want to report a bug or request a feature please raise an issue. To contribute code, please fork the repository and create a full request.

## Android app
You can get the [Android app](https://play.google.com/store/apps/details?id=com.keepassui) from the Google Play Store.


![keepassui_feature_image](https://github.com/npalumbo/keepassui/assets/1648970/0eb663cc-6fb3-44f3-86b1-a4b6daca4162)
