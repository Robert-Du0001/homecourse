/**
 * 课程数据
 */
export type CourseResource = {
  /**
   * 课程ID
   */
  id: number;
  /**
   * 课程标题
   */
  title: string;
  /**
   * 简介
   */
  description: string;
  /**
   * 封面路径
   */
  cover_path: string;
  /**
   * 创建日期
   */
  created_at: string;
};
