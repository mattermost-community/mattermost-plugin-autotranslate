import '@testing-library/jest-dom';
import React from 'react';
import {render, screen} from '@testing-library/react';

import MenuItem from './menu_item';

test('should render Translate menu when activated', async () => {
    render(<MenuItem/>);
    expect(screen.queryByText(/Translate/i)).toBeNull();

    render(<MenuItem activated={false}/>);
    expect(screen.queryByText(/Translate/i)).toBeNull();

    render(<MenuItem activated={true}/>);
    expect(screen.getByText(/Translate/i)).toBeInTheDocument();
});
