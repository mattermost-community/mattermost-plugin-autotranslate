import {combineReducers} from 'redux';

import {INFO_CHANGE, TRANSLATE_POST_SUCCESS} from './action_types';

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
    case TRANSLATE_POST_SUCCESS:
        const nextState = {};
        nextState[action.data.id] = action.data;

        return {...state, ...nextState};
    default:
        return state;
    }
};

export default combineReducers({
    translatedPosts,
    userInfo,
});
