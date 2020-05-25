import {combineReducers} from 'redux';

import {
    INFO_CHANGE,
    SAVE_TRANSLATED_POST,
    SAVE_TRANSLATION,
} from './action_types';

const userInfo = (state = {}, action) => {
    switch (action.type) {
    case INFO_CHANGE:
        return action.data;

    default:
        return state;
    }
};

const translatedPosts = (state = {}, action) => {
    switch (action.type) {
    case SAVE_TRANSLATED_POST: {
        const nextState = {};
        nextState[action.data.post_id] = action.data;

        return {...state, ...nextState};
    }
    default:
        return state;
    }
};

const translations = (state = {}, action) => {
    switch (action.type) {
    case SAVE_TRANSLATION: {
        const nextState = {};
        nextState[action.data.id] = action.data;

        return {...state, ...nextState};
    }
    default:
        return state;
    }
};

export default combineReducers({
    translatedPosts,
    translations,
    userInfo,
});
