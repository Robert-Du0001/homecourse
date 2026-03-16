/**
 * 课程分类数据
 */
export type CategoryResource = {
  /** 课程分类ID */
  id: number;
  /** 分类名 */
  name: string;
  /** 是否为默认分类（文件会默认扫描到此分类中） */
  is_default: boolean;
  /** 排序 */
  sort: number;
  /** 创建日期 */
  created_at: string;
};
