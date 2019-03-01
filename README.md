## Mattermost Autotranslation Plugin (beta) [![CircleCI](https://circleci.com/gh/mattermost/mattermost-plugin-autotranslate.svg?style=svg)](https://circleci.com/gh/mattermost/mattermost-plugin-autotranslate)

Autotranslation plugin for Mattermost.

Message autotranslation is powered by Amazon Translate which supports translation between English and any of the following languages: Arabic, Chinese (Simplified), Chinese (Traditional), Czech, French, German, Italian, Japanese, Portuguese, Russian, Spanish, and Turkish.

### Feature
* __Translate__ option available at dropdown menu of each regular post.
* __Slash commands__ to change user settings using `/autotranslate` slash command
    * __Check user info__ by issuing `/autotranslate info` to see current user setting
    * __Turn on/off__ translation by issuing `/autotranslate [on|off]`
    * __Change source language__ translation by initiating `/autotranslate source [language code]` (see language codes below)
    * __Change target language__ translation by initiating `/autotranslate target [language code]` (see language codes below)
* __Supported Languages and its codes__
    * __auto__ : Automatic language detection
    * __ar__ : Arabic
    * __zh__ : Chinese
    * __cs__ : Czech
    * __fr__ : French
    * __de__ : German
    * __en__ : English
    * __es__ : Spanish
    * __it__ : Italian
    * __ja__ : Japanese
    * __pt__ : Portuguese
    * __ru__ : Russian
    * __tr__ : Turkish

### Installation

__Requires Mattermost 5.4 or higher__
(Currently not compatible with Mattermost as there's minor change requested to plugin system.)

1. Install the plugin
    1. Download the latest version of the plugin from the GitHub releases page
    2. In Mattermost, go to the System Console -> Plugins -> Management
    3. Upload the plugin
2. Spin up Amazon Translate https://aws.amazon.com/translate/
3. In Mattermost, go to System Console -> Plugins -> Autotranslate
        * Fill in the AWS Access Key ID, Secret Access Key and Region
4. Enable the plugin
    * Go to System Console -> Plugins -> Management and click "Enable" underneath the Autotranslate plugin
6. Test it out
    * In Mattermost, run the slash command `/autotranslate on` and see if `Translate` option becomes available at dropdown menu of a post.

## Developing 

This plugin contains both a server and web app portion.

Use `make dist` to build distributions of the plugin that you can upload to a Mattermost server.

Use `make check-style` to check the style.

Use `make localdeploy` to deploy the plugin to your local server. You will need to restart the server to get the changes.
