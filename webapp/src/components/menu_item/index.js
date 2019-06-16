import {connect} from 'react-redux';

import {getUserInfo} from 'selectors';

import MenuItem from './menu_item';

const mapStateToProps = (state) => {
    const userInfo = getUserInfo(state);
    const activated = userInfo && userInfo.activated ? userInfo.activated : false;

    return {
        activated,
    };
};

export default connect(mapStateToProps)(MenuItem);
