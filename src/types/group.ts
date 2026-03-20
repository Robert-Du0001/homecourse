/**
 * 课程分组数据
 */
export type GroupResource = {
  /** 课程分组ID */
  id: number;
  /** 所属课程ID */
  course_id: number;
  /** 是否为默认分组（文件会默认扫描到此分组中） */
  is_default: boolean;
  /** 分组名 */
  name: string;
  /** 排序 */
  sort: number;
  /** 创建日期 */
  created_at: string;
};
