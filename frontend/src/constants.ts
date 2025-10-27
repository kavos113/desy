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
