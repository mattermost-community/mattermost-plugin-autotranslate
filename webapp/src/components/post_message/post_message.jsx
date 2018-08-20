import React from 'react';
import PropTypes from 'prop-types';

export default class PostMessage extends React.PureComponent {
    static propTypes = {
        activated: PropTypes.bool.isRequired,
        translation: PropTypes.object.isRequired,
        hide: PropTypes.func,
        onTextChange: PropTypes.func,
    }

    static defaultProps = {
        activated: false,
        translation: {},
    }

    componentDidUpdate(prevProps) {
        if (this.props.translation.translated_text !== prevProps.translation.translated_text) {
            this.props.onTextChange();
        }
    }

    handleOnClick = () => {
        this.props.hide(this.props.translation.post_id);
    }

    renderMessage(message) {
        return (
            <React.Fragment>
                <p>
                    <i className='icon fa fa-language'/>
                    {message}
                    <span
                        onClick={this.handleOnClick}
                        style={{cursor: 'pointer'}}
                    >
                        <i className='icon fa fa-chevron-circle-left'/>
                    </span>
                </p>
            </React.Fragment>
        );
    }

    render() {
        const {translation, activated} = this.props;

        if (!activated || !translation.show) {
            return null;
        }

        if (translation.error) {
            return this.renderMessage(
                <span style={{color: 'red'}}>{`  ${translation.error.status}: ${translation.error.text}  `}</span>
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
