import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import SearchBox from '../SearchBox';

describe('SearchBox', () => {
  it('未制御の入力値を更新できる', () => {
    const handleChange = vi.fn();
    render(<SearchBox placeholder="検索" onChange={handleChange} />);

    const input = screen.getByPlaceholderText('検索') as HTMLInputElement;
    fireEvent.change(input, { target: { value: 'アルゴリズム' } });

    expect(input.value).toBe('アルゴリズム');
    expect(handleChange).toHaveBeenCalledWith('アルゴリズム');
  });

  it('制御された入力値を保持する', () => {
    const handleChange = vi.fn();
    render(<SearchBox placeholder="検索" value="初期値" onChange={handleChange} />);

    const input = screen.getByPlaceholderText('検索') as HTMLInputElement;
    expect(input.value).toBe('初期値');

    fireEvent.change(input, { target: { value: '更新' } });
    expect(handleChange).toHaveBeenCalledWith('更新');
    expect(input.value).toBe('初期値');
  });
});
