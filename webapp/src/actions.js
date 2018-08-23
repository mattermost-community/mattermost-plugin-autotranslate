import {batchActions} from 'redux-batched-actions';

import {getPost} from 'mattermost-redux/selectors/entities/posts';
import {getCurrentUserId} from 'mattermost-redux/selectors/entities/users';
import {PostTypes} from 'mattermost-redux/action_types';

import {getTranslatedPosts, getUserInfo} from 'selectors';

import {
    INFO_CHANGE,
    TRANSLATE_POST_SUCCESS,
} from './action_types';

import Client from './clients';

export const postDropdownMenuAction = getTranslatedMessage;

export function getInfo() {
    return async (dispatch) => {
        let data;
        try {
            data = await Client.getInfo();
        } catch (error) {
            return {error};
        }

        dispatch({type: INFO_CHANGE, data});
        return {data};
    };
}

export function getTranslatedMessage(postId) {
    return async (dispatch, getState) => {
        const state = getState();

        const userInfo = getUserInfo(state);
        const {
            activated,
            source_language: source,
            target_language: target,
            user_id: userId,
        } = userInfo;

        if (!activated) {
            return {data: null};
        }

        const currentUserId = getCurrentUserId(state);
        if (currentUserId !== userId) {
            return {data: null};
        }

        const post = getPost(getState(), postId);
        if (!post) {
            return {data: null};
        }

        const translationKey = `${postId}${source}${target}${post.update_at}`;

        if (
            post &&
            post.translation &&
            post.translation.show &&
            post.translation.id === translationKey
        ) {
            return {data: post.translation};
        }

        const translationInStore = getTranslatedPosts(state)[translationKey];
        if (translationInStore && translationInStore.id) {
            dispatch({
                type: PostTypes.RECEIVED_POST,
                data: {...post, translation: {...translationInStore, show: true}},
            }, getState);
            return {data: translationInStore};
        }

        let data;
        try {
            data = await Client.getGo(postId, source, target);
        } catch (error) {
            const errorText = error.response && error.response.text ? error.response.text.split('\n')[0] : '';
            const text = errorText.replace(/[\n\t\r]/g, ' ');
            dispatch({
                type: PostTypes.RECEIVED_POST,
                data: {...post, translation: {errorMessage: text, show: true, post_id: postId}},
            }, getState);

            return {error};
        }

        const message = post.message + ' '; // ugly hack to re-render specific post which will cover show more/less
        dispatch(batchActions([
            {type: PostTypes.RECEIVED_POST, data: {...post, message, translation: {...data, show: true, error: ''}}},
            {type: TRANSLATE_POST_SUCCESS, data},
        ]), getState);

        return {data};
    };
}

export const hideTranslatedMessage = (postId) => (dispatch, getState) => {
    const post = getPost(getState(), postId);
    if (!post || !post.translation) {
        return {data: null};
    }

    const translation = post.translation;
    const hiddenTranslation = {...translation, show: false};

    dispatch({
        type: PostTypes.RECEIVED_POST,
        data: {...post, translation: hiddenTranslation},
    }, getState);

    return {data: translation};
};

export const websocketInfoChange = (message) => (dispatch) => dispatch({
    type: INFO_CHANGE,
    data: message.data,
});
