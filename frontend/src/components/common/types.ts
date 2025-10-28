export interface MenuRecord {
  [key: string]: Menu;
}

export type Menu = string[] | MenuRecord;
