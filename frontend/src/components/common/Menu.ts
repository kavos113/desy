export type MenuLeaf = string[];

export interface MenuNode {
  [key: string]: MenuNode | MenuLeaf;
}

export type Menu = MenuLeaf | MenuNode;

export function isMenuNode(value: unknown): value is MenuNode {
  return Boolean(value) && typeof value === "object" && !Array.isArray(value);
}

export function isMenuLeaf(value: unknown): value is MenuLeaf {
  return Array.isArray(value);
}
