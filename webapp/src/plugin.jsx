import React from 'react';

import {getPost} from 'mattermost-redux/selectors/entities/posts';

import PostMessageAttachment from './components/post_message_attachment';
import TranslateMenuItem from './components/translate_menu_item';

import PluginId from './plugin_id';

import {
    getTranslatedMessage,
    getInfo,
    websocketInfoChange,
} from './actions';
import reducer from './reducer';
import {getUserInfo} from './selectors';

export default class AWSTranslatePlugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry, store) {
        registry.registerReducer(reducer);

        // Immediately fetch the current plugin status.
        store.dispatch(getInfo());

        registry.registerPostMessageAttachmentComponent(PostMessageAttachment);
        registry.registerPostDropdownMenuAction(
            <TranslateMenuItem/>,
            (postId) => store.dispatch(getTranslatedMessage(postId)),
            (postId) => {
                const state = store.getState();
                const post = getPost(state, postId);
                const userInfo = getUserInfo(state);
                return post && post.type === '' && userInfo && userInfo.activated;
            },
        );

        registry.registerWebSocketEventHandler(
            'custom_' + PluginId + '_info_change',
            (message) => {
                store.dispatch(websocketInfoChange(message));
            },
        );

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
