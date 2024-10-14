import React from 'react';
import { render, screen } from '@testing-library/react';
import MyComponent from './MyComponent';

test('renders the name prop correctly', () => {
  render(<MyComponent name='John' />);
  const headingElement = screen.getByText(/Hello, John!/i);
  expect(headingElement).toBeInTheDocument();
});
