import { act, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import SearchBoxes from '../SearchBoxes';

describe('SearchBoxes', () => {
  it('検索ボックスの入力をコールバックへ伝播する', async () => {
    const handleChange = vi.fn();
    const user = userEvent.setup();
    render(<SearchBoxes onChangeSearchBox={handleChange} />);

    expect(handleChange).not.toHaveBeenCalled();

    const titleInput = screen.getByPlaceholderText('講義名');
    await act(async () => {
      await user.type(titleInput, '法学');
    });

    expect(handleChange).toHaveBeenCalled();
  expect(handleChange).toHaveBeenLastCalledWith('法学', '', '');

    const lecturerInput = screen.getByPlaceholderText('教員名');
    await act(async () => {
      await user.type(lecturerInput, '佐藤');
    });

    expect(handleChange).toHaveBeenLastCalledWith('法学', '佐藤', '');

    const roomInput = screen.getByPlaceholderText('講義室名');
    await act(async () => {
      await user.type(roomInput, '101');
    });

    expect(handleChange).toHaveBeenLastCalledWith('法学', '佐藤', '101');
  });
});
