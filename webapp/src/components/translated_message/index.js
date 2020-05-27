import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getUserInfo, getTranslatedPosts} from 'selectors';
import {hideTranslatedMessage} from 'actions';

import TranslatedMessage from './translated_message';

const mapStateToProps = (state, ownProps) => {
    const userInfo = getUserInfo(state);
    const activated = userInfo && userInfo.activated ? userInfo.activated : false;

    return {
        activated,
        translation: getTranslatedPosts(state)[ownProps.postId],
    };
};

const mapDispatchToProps = (dispatch) => bindActionCreators({
    hideTranslatedMessage,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(TranslatedMessage);
