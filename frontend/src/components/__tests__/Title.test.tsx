import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import Title from '../Title';

describe('Title', () => {
  it('Pocket Syllabus の見出しを表示する', () => {
    render(<Title />);
    expect(screen.getByText('Pocket Syllabus')).toBeInTheDocument();
  });
});
