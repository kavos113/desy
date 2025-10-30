import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import ListHeaderItem from '../ListHeaderItem';

describe('ListHeaderItem', () => {
  it('クリックした列のキーでソートイベントを発火する', () => {
    const handleSort = vi.fn();
    render(<ListHeaderItem onSort={handleSort} />);

    fireEvent.click(screen.getByText('コード'));

    expect(handleSort).toHaveBeenCalledWith('code');
  });
});
