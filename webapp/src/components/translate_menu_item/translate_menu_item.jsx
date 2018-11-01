import PropTypes from 'prop-types';

const TranslateMenuItem = ({activated}) => {
    if (!activated) {
        return null;
    }

    return 'Translate';
};

TranslateMenuItem.propTypes = {
    activated: PropTypes.bool,
};

export default TranslateMenuItem;
