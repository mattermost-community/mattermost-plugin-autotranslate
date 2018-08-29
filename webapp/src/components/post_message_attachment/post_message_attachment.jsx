import React from 'react';
import PropTypes from 'prop-types';

export default class PostMessageAttachment extends React.PureComponent {
    static propTypes = {
        activated: PropTypes.bool.isRequired,
        translation: PropTypes.object.isRequired,
        hide: PropTypes.func,
        onHeightChange: PropTypes.func,
    }

    static defaultProps = {
        activated: false,
        translation: {},
    }

    componentDidUpdate(prevProps) {
        if (this.props.translation.translated_text !== prevProps.translation.translated_text) {
            this.props.onHeightChange(1);
        }
    }

    handleOnClick = () => {
        this.props.hide(this.props.translation.post_id);
        this.props.onHeightChange(1);
    }

    renderMessage(message) {
        return (
            <React.Fragment>
                <p>
                    <i className='icon fa fa-language'/>
                    {message}
                    <a onClick={this.handleOnClick}>{'(close)'}</a>
                </p>
            </React.Fragment>
        );
    }

    render() {
        const {translation, activated} = this.props;

        if (!activated || !translation.show) {
            return null;
        }

        if (translation.errorMessage) {
            return this.renderMessage(
                <span style={{color: 'red'}}>{`  ${translation.errorMessage}  `}</span>
            );
        }

        return this.renderMessage(
            <React.Fragment>
                <span>{'  See translation:\n'}</span>
                <span>{`${translation.translated_text}  `}</span>
            </React.Fragment>
        );
    }
}
