import React from 'react';
import PropTypes from 'prop-types';

export default class ErrorBoundary extends React.PureComponent {
    constructor(props) {
        super(props);
        this.state = {hasError: false};
    }

    static propTypes = {
        children: PropTypes.node,
    }

    static getDerivedStateFromError() {
        return {hasError: true};
    }

    render() {
        if (this.state.hasError) {
            return null;
        }

        return this.props.children;
    }
}
