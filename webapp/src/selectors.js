import PluginId from './plugin_id';

const getPluginState = (state) => state['plugins-' + PluginId] || {};

export const getUserInfo = (state) => getPluginState(state).userInfo;
export const getTranslatedPosts = (state) => getPluginState(state).translatedPosts;
export const getTranslations = (state) => getPluginState(state).translations;
