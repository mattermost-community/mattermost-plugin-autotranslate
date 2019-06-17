import ErrorBoundary from './error_boundary';
import TranslatedMessage from './translated_message';

const PostMessageAttachment = () => {
    return (
        <ErrorBoundary>
            <TranslatedMessage/>
        </ErrorBoundary>
    );
};

export default PostMessageAttachment;
