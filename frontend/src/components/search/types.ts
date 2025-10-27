export type Day = "月" | "火" | "水" | "木" | "金";
export type Period = "1" | "2" | "3" | "4" | "5";

export type SearchComboBox = "university" | "department" | "year";
export type SearchSearchBox = "title" | "lecturer";
export type SearchCheckBox = "grade" | "quarter";
export type SearchConditionsType =
  | SearchComboBox
  | SearchSearchBox
  | SearchCheckBox;

export type SearchTimetableCell = {
  day: Day;
  period: Period;
};

export const DAY_OPTIONS: Day[] = ["月", "火", "水", "木", "金"];
export const PERIOD_OPTIONS: Period[] = ["1", "2", "3", "4", "5"];
export const QUARTER_OPTIONS = ["1Q", "2Q", "3Q", "4Q"];
export const GRADE_LABELS = [
  "学士1年",
  "学士2年",
  "学士3年",
  "修士1年",
  "修士2年",
  "博士課程",
];
