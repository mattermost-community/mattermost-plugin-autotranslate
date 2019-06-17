import ErrorBoundary from './error_boundary';
import MenuItem from './menu_item';

const TranslateMenuItem = () => {
    return (
        <ErrorBoundary>
            <MenuItem/>
        </ErrorBoundary>
    );
};

export default TranslateMenuItem;
