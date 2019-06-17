import PropTypes from 'prop-types';

const MenuItem = ({activated}) => {
    if (!activated) {
        return null;
    }

    return 'Translate';
};

MenuItem.propTypes = {
    activated: PropTypes.bool,
};

export default MenuItem;
