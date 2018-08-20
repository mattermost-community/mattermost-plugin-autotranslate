import request from 'superagent';

import PluginId from './plugin_id';
import {buildQueryString} from './utils';

class ClientClass {
    constructor() {
        this.url = `/plugins/${PluginId}/api`;
    }

    getGo = async (postId, source, target) => {
        return this.doGet(this.url + '/go' + buildQueryString({post_id: postId, source, target}));
    }

    getInfo = async () => {
        return this.doGet(`${this.url}/get_info`);
    }

    postInfo = async (info) => {
        return this.doPost(`${this.url}/set_info`, info);
    }

    doGet = async (url, headers = {}) => {
        headers['X-Requested-With'] = 'XMLHttpRequest';

        try {
            const response = await request.
                get(url).
                set(headers).
                accept('application/json');

            return response.body;
        } catch (err) {
            throw err;
        }
    }

    doPost = async (url, body, headers = {}) => {
        headers['X-Requested-With'] = 'XMLHttpRequest';

        try {
            const response = await request.
                post(url).
                send(body).
                set(headers).
                type('application/json').
                accept('application/json');

            return response.body;
        } catch (err) {
            throw err;
        }
    }
}

const Client = new ClientClass();

export default Client;
