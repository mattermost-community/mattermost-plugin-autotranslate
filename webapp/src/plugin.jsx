import React from 'react';

import PluginId from './plugin_id';

import PostMessage from './components/post_message';
import TranslateMenuItem from './components/translate_menu_item';

import {
    postDropdownMenuAction,
    getInfo,
    websocketInfoChange,
} from './actions';
import reducer from './reducer';

export default class AWSTranslatePlugin {
    initialize(registry, store) {
        registry.registerPostMessageComponent(PostMessage);
        registry.registerPostDropdownMenuAction(
            <TranslateMenuItem/>,
            (post) => store.dispatch(postDropdownMenuAction(post)),
            (post) => post && post.type === '',
        );

        registry.registerWebSocketEventHandler(
            'custom_' + PluginId + '_info_change',
            (message) => {
                store.dispatch(websocketInfoChange(message));
            },
        );

        registry.registerReducer(reducer);

        // Immediately fetch the current plugin status.
        store.dispatch(getInfo());

        // Fetch the current status whenever we recover an internet connection.
        registry.registerReconnectHandler(() => {
            store.dispatch(getInfo());
        });
    }

    uninitialize() {
        //eslint-disable-next-line no-console
        console.log(PluginId + '::uninitialize()');
    }
}
