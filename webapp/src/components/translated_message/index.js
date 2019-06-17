import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';

import {getPost} from 'mattermost-redux/selectors/entities/posts';

import {getUserInfo} from 'selectors';
import {hideTranslatedMessage} from 'actions';

import TranslatedMessage from './translated_message';

const mapStateToProps = (state, ownProps) => {
    const userInfo = getUserInfo(state);
    const activated = userInfo && userInfo.activated ? userInfo.activated : false;
    const post = getPost(state, ownProps.postId);
    const translation = post && post.translation ? post.translation : null;
    return {
        activated,
        translation,
    };
};

const mapDispatchToProps = (dispatch) => bindActionCreators({
    hide: hideTranslatedMessage,
}, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(TranslatedMessage);
