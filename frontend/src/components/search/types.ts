import type { Day, Period } from '../../constants';

export type SearchComboBox = 'university' | 'department' | 'year';
export type SearchSearchBox = 'title' | 'lecturer';
export type SearchCheckBox = 'grade' | 'quarter';
export type SearchConditionKey = SearchComboBox | SearchSearchBox | SearchCheckBox;

export interface SearchTimetableSelection {
  day: Day;
  period: Period;
}

export interface SearchState {
  university: string[];
  department: string[];
  year: string[];
  title: string[];
  lecturer: string[];
  grade: string[];
  quarter: string[];
  timetable: SearchTimetableSelection[];
}
