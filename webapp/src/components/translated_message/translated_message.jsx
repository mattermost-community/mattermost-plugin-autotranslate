import React from 'react';
import PropTypes from 'prop-types';

export default class TranslatedMessage extends React.PureComponent {
    static propTypes = {
        activated: PropTypes.bool.isRequired,
        translation: PropTypes.object,
        hideTranslatedMessage: PropTypes.func.isRequired,
        onHeightChange: PropTypes.func,
    }

    static defaultProps = {
        activated: false,
    }

    componentDidUpdate(prevProps) {
        if (this.props.translation &&
            prevProps.translation &&
            this.props.translation.translated_text !== prevProps.translation.translated_text
        ) {
            this.props.onHeightChange(1);
        }
    }

    handleCloseMessage = () => {
        this.props.hideTranslatedMessage(this.props.translation.post_id);
        this.props.onHeightChange(1);
    }

    renderMessage(message) {
        return (
            <React.Fragment>
                <p>
                    <i className='icon fa fa-language'/>
                    {message}
                    <a onClick={this.handleCloseMessage}>{'(close)'}</a>
                </p>
            </React.Fragment>
        );
    }

    render() {
        const {translation, activated} = this.props;

        if (!activated || !translation || !translation.show) {
            return null;
        }

        if (translation.errorMessage) {
            return this.renderMessage(
                <span style={{color: 'red'}}>{`  ${translation.errorMessage}  `}</span>,
            );
        }

        return this.renderMessage(
            <React.Fragment>
                <span>{'  See translation:\n'}</span>
                <span>{`${translation.translated_text}  `}</span>
            </React.Fragment>,
        );
    }
}
