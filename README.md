## Mattermost Autotranslation Plugin (Beta)

[![Build Status](https://img.shields.io/circleci/project/github/mattermost/mattermost-plugin-autotranslate/master)](https://circleci.com/gh/mattermost/mattermost-plugin-autotranslate)
[![Code Coverage](https://img.shields.io/codecov/c/github/mattermost/mattermost-plugin-autotranslate/master)](https://codecov.io/gh/mattermost/mattermost-plugin-autotranslate)
[![Release](https://img.shields.io/github/v/release/mattermost/mattermost-plugin-autotranslate)](https://github.com/mattermost/mattermost-plugin-autotranslate/releases/latest)
[![HW](https://img.shields.io/github/issues/mattermost/mattermost-plugin-autotranslate/Up%20For%20Grabs?color=dark%20green&label=Help%20Wanted)](https://github.com/mattermost/mattermost-plugin-autotranslate/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22Up+For+Grabs%22+label%3A%22Help+Wanted%22)

**Maintainer:** [@saturninoabril](https://github.com/saturninoabril)

### Autotranslation plugin for Mattermost.

Message autotranslation is powered by Amazon Translate which is a text translation service that uses advanced machine learning technologies to provide high-quality translation on demand. Amazon Translate can translate text between the languages listed in its [website](https://docs.aws.amazon.com/translate/latest/dg/what-is.html).

### Feature
* __Translate__ option available at dropdown menu of each regular post.
* __Slash commands__ to change user settings using `/autotranslate` slash command
    * __Check user info__ by issuing `/autotranslate info` to see current user setting
    * __Turn on/off__ translation by issuing `/autotranslate [on|off]`
    * __Change source language__ translation by initiating `/autotranslate source [language code]`
    * __Change target language__ translation by initiating `/autotranslate target [language code]`
* __Supported Languages and its codes__ can be found at [Amazon Translate website](https://docs.aws.amazon.com/translate/latest/dg/what-is.html). 

### Installation

__Requires Mattermost 5.22 or higher__

1. Install the plugin
    1. Download the latest version of the plugin from the GitHub releases page
    2. In Mattermost, go to the System Console -> Plugins -> Management
    3. Upload the plugin
2. Spin up Amazon Translate https://aws.amazon.com/translate/
3. In Mattermost, go to System Console -> Plugins -> Autotranslate
        * Fill in the AWS Access Key ID, Secret Access Key and Region
4. Enable the plugin
    * Go to System Console -> Plugins -> Management and click "Enable" underneath the Autotranslate plugin
5. Test it out
    * In Mattermost, run the slash command `/autotranslate on` and see if `Translate` option becomes available at dropdown menu of a post.

## Developing 

This plugin contains both a server and web app portion.

Use `make dist` to build distributions of the plugin that you can upload to a Mattermost server.

Use `make check-style` to check the style.

Use `make localdeploy` to deploy the plugin to your local server. You will need to restart the server to get the changes.
