import React from 'react';
import PropTypes from 'prop-types';

const MenuItem = ({activated}) => {
    if (!activated) {
        return null;
    }

    return (
        <button
            className='style--none'
            role='presentation'
        >
            <span className='MenuItem__icon'>
                <i className='icon fa fa-language'/>
            </span>
            <span>{'Translate'}</span>
        </button>
    );
};

MenuItem.propTypes = {
    activated: PropTypes.bool,
};

export default MenuItem;
