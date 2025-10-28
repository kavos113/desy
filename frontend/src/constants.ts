import type { Menu } from "./common/Menu";

const RAW_DEPARTMENTS = [
  "理学院",
  "工学院",
  "物質理工学院",
  "情報理工学院",
  "生命理工学院",
  "環境・社会理工学院",
  "工学院，物質理工学院，環境・社会理工学院共通科目",
  "教養科目群",
  "数学系",
  "物理学系",
  "化学系",
  "地球惑星科学系",
  "数学コース",
  "物理学コース",
  "化学コース",
  "エネルギーコース",
  "エネルギー・情報コース",
  "地球惑星科学コース",
  "地球生命コース",
  "機械系",
  "システム制御系",
  "電気電子系",
  "情報通信系",
  "経営工学系",
  "機械コース",
  "エンジニアリングデザインコース",
  "ライフエンジニアリングコース",
  "原子核工学コース",
  "システム制御コース",
  "電気電子コース",
  "情報通信コース",
  "経営工学コース",
  "材料系",
  "応用化学系",
  "材料コース",
  "応用化学コース",
  "数理・計算科学系",
  "情報工学系",
  "数理・計算科学コース",
  "情報工学コース",
  "知能情報コース",
  "生命理工学系",
  "生命理工学コース",
  "建築学系",
  "土木・環境工学系",
  "融合理工学系",
  "建築学コース",
  "土木工学コース",
  "融合理工学コース",
  "都市・環境学コース",
  "地球環境共創コース",
  "社会・人間科学コース",
  "イノベーション科学コース",
  "技術経営専門職学位課程",
  "文系教養科目",
  "英語科目",
  "第二外国語科目",
  "日本語・日本文化科目",
  "教職科目",
  "アントレプレナーシップ科目",
  "広域教養科目",
  "理工系教養科目",
  "キャリア科目",
];

const collator = new Intl.Collator("ja-JP");

export const DEPARTMENT_OPTIONS = Array.from(new Set(RAW_DEPARTMENTS)).sort(
  (a, b) => collator.compare(a, b)
);

export const LEVEL_OPTIONS = [
  { value: 1, label: "学士1年" },
  { value: 2, label: "学士2年" },
  { value: 3, label: "学士3年" },
  { value: 4, label: "修士1年" },
  { value: 5, label: "修士2年" },
  { value: 6, label: "博士課程" },
];

const CURRENT_YEAR = new Date().getFullYear();
const YEAR_RANGE = 4;

export const YEAR_OPTIONS = Array.from(
  { length: YEAR_RANGE },
  (_, index) => CURRENT_YEAR - index
);

export const KEYWORD_SEPARATOR = /[\s,\u3001\u3002、;；]+/g;

export function parseKeywordInput(value: string): string[] {
  return value
    .split(KEYWORD_SEPARATOR)
    .map((item) => item.trim())
    .filter((item) => item.length > 0);
}

export const LECTURE_TYPE_LABELS: Record<string, string> = {
  offline: "対面",
  live: "ライブ",
  hyflex: "ハイフレックス",
  ondemand: "オンデマンド",
  other: "その他",
};

export const LEVEL_LABELS = LEVEL_OPTIONS.reduce<Record<number, string>>(
  (acc, option) => {
    acc[option.value] = option.label;
    return acc;
  },
  {}
);

export type Day = "月" | "火" | "水" | "木" | "金";

export const DAYS: Day[] = ["月", "火", "水", "木", "金"];

export type Period = "1" | "2" | "3" | "4" | "5";

export const PERIODS: Period[] = ["1", "2", "3", "4", "5"];

export const GRADE_LABELS = LEVEL_OPTIONS.map((option) => option.label);

export const QUARTER_LABELS = ["1Q", "2Q", "3Q", "4Q"];

export const UNIVERSITIES_MENU: Menu = {
  大学を選択: ["東京工業大学", "一橋大学"],
};

export const DEPARTMENTS_MENU: Menu = {
  開講元を選択: {
    東京工業大学: {
      学士課程: {
        理学院: ["数学系", "物理学系", "化学系", "地球惑星科学系"],
        工学院: [
          "機械系",
          "システム制御系",
          "電気電子系",
          "情報通信系",
          "経営工学系",
        ],
        物質理工学院: ["材料系", "応用化学系"],
        情報理工学院: ["数理・計算科学系", "情報工学系"],
        生命理工学院: ["生命理工学系"],
        "環境・社会理工学院": ["建築学系", "土木・環境工学系", "融合理工学系"],
        "工学院，物質理工学院，環境・社会理工学院共通科目": [
          "工学院，物質理工学院，環境・社会理工学院共通科目",
        ],
        教養科目群: [
          "文系教養科目",
          "英語科目",
          "第二外国語科目",
          "日本語・日本文化科目",
          "教職科目",
          "アントレプレナーシップ科目",
          "広域教養科目",
          "理工系教養科目",
        ],
      },
      大学院課程: {
        理学院: [
          "数学コース",
          "物理学コース",
          "化学コース",
          "エネルギーコース",
          "エネルギー・情報コース",
          "地球惑星科学コース",
          "地球生命コース",
        ],
        工学院: [
          "機械コース",
          "エネルギーコース",
          "エネルギー・情報コース",
          "エンジニアリングデザインコース",
          "ライフエンジニアリングコース",
          "原子核工学コース",
          "システム制御コース",
          "電気電子コース",
          "情報通信コース",
          "経営工学コース",
        ],
        物質理工学院: [
          "材料コース",
          "応用化学コース",
          "エネルギーコース",
          "エネルギー・情報コース",
          "ライフエンジニアリングコース",
          "原子核工学コース",
          "地球生命コース",
        ],
        情報理工学院: [
          "数理・計算科学コース",
          "情報工学コース",
          "知能情報コース",
          "エネルギー・情報コース",
          "ライフエンジニアリングコース",
        ],
        生命理工学院: [
          "生命理工学コース",
          "ライフエンジニアリングコース",
          "地球生命コース",
        ],
        "環境・社会理工学院": [
          "建築学コース",
          "土木工学コース",
          "融合理工学コース",
          "エンジニアリングデザインコース",
          "都市・環境学コース",
          "地球環境共創コース",
          "エネルギーコース",
          "エネルギー・情報コース",
          "原子核工学コース",
          "社会・人間科学コース",
          "イノベーション科学コース",
          "技術経営専門職学位課程",
        ],
        教養科目群: [
          "文系教養科目",
          "英語科目",
          "第二外国語科目",
          "日本語・日本文化科目",
          "教職科目",
          "アントレプレナーシップ科目",
          "広域教養科目",
          "キャリア科目",
        ],
      },
    },
  },
};

export const MOBILE_DEPARTMENTS_MENU: Menu = {
  開講元を選択: [
    "理学院",
    "工学院",
    "物質理工学院",
    "情報理工学院",
    "生命理工学院",
    "環境・社会理工学院",
    "工学院，物質理工学院，環境・社会理工学院共通科目",
    "教養科目群",
    "数学系",
    "物理学系",
    "化学系",
    "地球惑星科学系",
    "数学コース",
    "物理学コース",
    "化学コース",
    "エネルギーコース",
    "エネルギー・情報コース",
    "地球惑星科学コース",
    "地球生命コース",
    "機械系",
    "システム制御系",
    "電気電子系",
    "情報通信系",
    "経営工学系",
    "機械コース",
    "エンジニアリングデザインコース",
    "ライフエンジニアリングコース",
    "原子核工学コース",
    "システム制御コース",
    "電気電子コース",
    "情報通信コース",
    "経営工学コース",
    "材料系",
    "応用化学系",
    "材料コース",
    "応用化学コース",
    "数理・計算科学系",
    "情報工学系",
    "数理・計算科学コース",
    "情報工学コース",
    "知能情報コース",
    "生命理工学系",
    "生命理工学コース",
    "環境・社会理工学院",
    "建築学系",
    "土木・環境工学系",
    "融合理工学系",
    "建築学コース",
    "土木工学コース",
    "融合理工学コース",
    "都市・環境学コース",
    "地球環境共創コース",
    "社会・人間科学コース",
    "イノベーション科学コース",
    "技術経営専門職学位課程",
    "文系教養科目",
    "英語科目",
    "第二外国語科目",
    "日本語・日本文化科目",
    "教職科目",
    "アントレプレナーシップ科目",
    "広域教養科目",
    "理工系教養科目",
    "キャリア科目",
  ],
};

export const YEARS_MENU: Menu = {
  年度を選択: Array.from(
    { length: YEAR_RANGE },
    (_, index) => `${CURRENT_YEAR - index}年度`
  ),
};

const DAY_TO_DOMAIN: Record<Day, string> = {
  月: "monday",
  火: "tuesday",
  水: "wednesday",
  木: "thursday",
  金: "friday",
};

export function dayToDomain(day: Day): string {
  return DAY_TO_DOMAIN[day];
}

export function periodToNumber(period: Period): number {
  return Number.parseInt(period, 10);
}

export function gradeLabelToLevel(label: string): number | undefined {
  const option = LEVEL_OPTIONS.find((item) => item.label === label);
  return option?.value;
}

const QUARTER_TO_SEMESTER: Record<string, string> = {
  "1Q": "spring",
  "2Q": "summer",
  "3Q": "fall",
  "4Q": "winter",
};

export function quarterToSemesters(quarter: string): string[] {
  const semester = QUARTER_TO_SEMESTER[quarter];
  return semester ? [semester] : [];
}

export function parseYearLabel(label: string): number | undefined {
  const match = label.match(/(\d{4})/);
  if (!match) {
    return undefined;
  }
  return Number.parseInt(match[1], 10);
}
