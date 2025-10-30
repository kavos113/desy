import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import ComboBox from '../ComboBox';
import { Menu } from '../Menu';

const SIMPLE_MENU: Menu = {
  Parent: ['ChildA', 'ChildB']
};

describe('ComboBox', () => {
  it('トップレベルの項目選択をトグルできる', () => {
    const handleSelect = vi.fn();
    render(<ComboBox items={SIMPLE_MENU} onSelectItem={handleSelect} />);

    const topButton = screen.getByRole('button', { name: 'Parent' });
    fireEvent.click(topButton);
    expect(handleSelect).toHaveBeenLastCalledWith(['Parent']);

    fireEvent.click(topButton);
    expect(handleSelect).toHaveBeenLastCalledWith([]);
  });

  it('サブ項目を複数選択できる', () => {
    const handleSelect = vi.fn();
    render(<ComboBox items={SIMPLE_MENU} onSelectItem={handleSelect} />);

    const parentButton = screen.getByRole('button', { name: 'Parent' });
    const parentItem = parentButton.closest('li');
    if (!parentItem) {
      throw new Error('親項目のリスト要素が見つかりません');
    }

    fireEvent.mouseEnter(parentItem);

    const childA = screen.getByText('ChildA');
    const childB = screen.getByText('ChildB');

    fireEvent.click(childA);
    expect(handleSelect).toHaveBeenLastCalledWith(['ChildA']);

    fireEvent.click(childB);
    expect(handleSelect).toHaveBeenLastCalledWith(['ChildA', 'ChildB']);
  });
});
