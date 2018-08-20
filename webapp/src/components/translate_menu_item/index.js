import {connect} from 'react-redux';

import {getUserInfo} from 'selectors';

import TranslateMenuItem from './translate_menu_item';

const mapStateToProps = (state) => {
    const userInfo = getUserInfo(state) || {};

    return {
        activated: userInfo.activated,
    };
};

export default connect(mapStateToProps)(TranslateMenuItem);
