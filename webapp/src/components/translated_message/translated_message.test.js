import '@testing-library/jest-dom';
import React from 'react';
import {fireEvent, render, screen} from '@testing-library/react';

import TranslatedMessage from './translated_message';

test('should not render any message', async () => {
    const hideTranslatedMessage = jest.fn();
    const {container, rerender} = render(
        <TranslatedMessage
            activated={false}
            hideTranslatedMessage={hideTranslatedMessage}
        />,
    );
    expect(container).toMatchInlineSnapshot('<div />');

    rerender(
        <TranslatedMessage
            translation={{show: false}}
            activated={false}
            hideTranslatedMessage={hideTranslatedMessage}
        />,
    );
    expect(container).toMatchSnapshot();
});

test('should render error message', async () => {
    const hideTranslatedMessage = jest.fn();
    const onHeightChange = jest.fn();
    const translation = {
        post_id: 'post_id_1',
        show: true,
        errorMessage: 'Failed to get translation',
    };

    const {container} = render(
        <TranslatedMessage
            translation={translation}
            activated={true}
            hideTranslatedMessage={hideTranslatedMessage}
            onHeightChange={onHeightChange}
        />,
    );
    expect(container).toMatchSnapshot();
    expect(screen.getByText(translation.errorMessage)).toBeInTheDocument();

    fireEvent.click(screen.getByText(/close/i));
    expect(hideTranslatedMessage).toHaveBeenCalledTimes(1);
    expect(hideTranslatedMessage).toHaveBeenCalledWith(translation.post_id);
    expect(onHeightChange).toHaveBeenCalledTimes(1);
    expect(onHeightChange).toHaveBeenCalledWith(1);
});

test('should render translated message', async () => {
    const hideTranslatedMessage = jest.fn();
    const onHeightChange = jest.fn();
    const translation = {
        post_id: 'post_id_2',
        show: true,
        translated_text: 'Hello world',
    };

    const {container} = render(
        <TranslatedMessage
            translation={translation}
            activated={true}
            hideTranslatedMessage={hideTranslatedMessage}
            onHeightChange={onHeightChange}
        />,
    );
    expect(container).toMatchSnapshot();
    expect(screen.getByText(translation.translated_text)).toBeInTheDocument();
    expect(screen.getByText(/See translation/i)).toBeInTheDocument();

    fireEvent.click(screen.getByText(/close/i));
    expect(hideTranslatedMessage).toHaveBeenCalledTimes(1);
    expect(hideTranslatedMessage).toHaveBeenCalledWith(translation.post_id);
    expect(onHeightChange).toHaveBeenCalledTimes(1);
    expect(onHeightChange).toHaveBeenCalledWith(1);
});
