/**
 * 课程分组数据
 */
export type GroupResource = {
  /** 课程分组ID */
  id: number;
  /** 分类ID */
  category_id: number;
  /** 分组名 */
  name: string;
  /** 排序 */
  sort: number;
  /** 创建日期 */
  created_at: string;
};
