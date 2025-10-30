import { describe, expect, it } from 'vitest';
import { formatTimetables } from '../utils';

describe('formatTimetables', () => {
  it('連続する同曜日・同教室の時限をまとめて表示する', () => {
    const result = formatTimetables([
      { DayOfWeek: 'monday', Period: 1, Room: { Name: 'W1' } } as any,
      { DayOfWeek: 'monday', Period: 2, Room: { Name: 'W1' } } as any,
      { DayOfWeek: 'monday', Period: 3, Room: { Name: 'W1' } } as any
    ]);

    expect(result).toBe('月1-3(W1)');
  });

  it('オプションで教室表示を省略できる', () => {
    const result = formatTimetables(
      [
        { DayOfWeek: 'friday', Period: 5, Room: { Name: 'M110' } } as any,
        { DayOfWeek: 'friday', Period: 6, Room: { Name: 'M110' } } as any
      ],
      { includeRoom: false }
    );

    expect(result).toBe('金5-6');
  });

  it('同じ曜日でも教室が異なる場合は別々に表示する', () => {
    const result = formatTimetables([
      { DayOfWeek: 'tuesday', Period: 3, Room: { Name: 'W1' } } as any,
      { DayOfWeek: 'tuesday', Period: 4, Room: { Name: 'W2' } } as any
    ]);

    expect(result).toBe('火3(W1), 火4(W2)');
  });

  it('連続していない時限は個別に表示する', () => {
    const result = formatTimetables([
      { DayOfWeek: 'wednesday', Period: 2, Room: { Name: 'W1' } } as any,
      { DayOfWeek: 'wednesday', Period: 4, Room: { Name: 'W1' } } as any
    ]);

    expect(result).toBe('水2(W1), 水4(W1)');
  });
});
