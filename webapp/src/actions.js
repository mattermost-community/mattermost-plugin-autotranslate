import {getPost} from 'mattermost-redux/selectors/entities/posts';
import {getCurrentUserId} from 'mattermost-redux/selectors/entities/users';

import {
    getTranslatedPosts,
    getTranslations,
    getUserInfo,
} from 'selectors';

import {
    INFO_CHANGE,
    SAVE_TRANSLATED_POST,
    SAVE_TRANSLATION,
} from './action_types';

import Client from './clients';

export const getInfo = () => {
    return async (dispatch) => {
        try {
            const data = await Client.getInfo();
            dispatch({type: INFO_CHANGE, data});

            return {data};
        } catch (error) {
            return {error};
        }
    };
};

export const getTranslatedMessage = (postId) => {
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

        const translationId = `${postId}${source}${target}${post.update_at}`;
        const translation = getTranslations(state)[translationId];
        if (translation) {
            dispatch(saveTranslatedPost({...translation, show: true}));
            return {success: true};
        }

        let result;
        try {
            result = await Client.getGo(postId, source, target);
        } catch (error) {
            const errorText = error.response && error.response.text ? error.response.text.split('\n')[0] : '';
            const text = errorText.replace(/[\n\t\r]/g, ' ');
            const errorData = {errorMessage: text, show: true, post_id: postId};

            dispatch(saveTranslatedPost(errorData));
            return {error: true};
        }

        dispatch(saveTranslatedPost({...result, show: true}));
        dispatch(saveTranslation(result));

        return {success: true};
    };
};

export const saveTranslatedPost = (data) => {
    return (dispatch) => {
        dispatch({type: SAVE_TRANSLATED_POST, data});
    };
};

export const saveTranslation = (data) => {
    return (dispatch) => {
        dispatch({type: SAVE_TRANSLATION, data});
    };
};

export const hideTranslatedMessage = (postId) => {
    return (dispatch, getState) => {
        const translatedPost = getTranslatedPosts(getState())[postId];
        if (!translatedPost) {
            return;
        }

        const hiddenTranslation = {...translatedPost, show: false};
        dispatch({type: SAVE_TRANSLATED_POST, data: hiddenTranslation}, getState);
    };
};

export const websocketInfoChange = (message) => {
    return (dispatch) => {
        dispatch({type: INFO_CHANGE, data: message.data});
    };
};
