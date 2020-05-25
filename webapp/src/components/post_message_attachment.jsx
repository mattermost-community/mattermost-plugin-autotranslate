import ErrorBoundary from './error_boundary';
import TranslatedMessage from './translated_message';

const PostMessageAttachment = (props = {}) => {
    return (
        <ErrorBoundary>
            <TranslatedMessage {...props}/>
        </ErrorBoundary>
    );
};

export default PostMessageAttachment;
