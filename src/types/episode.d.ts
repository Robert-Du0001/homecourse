import type { CourseResource } from '@/types/course';

/**
 * 剧集数据
 */
export type EpisodesResource = {
  /**
   * 剧集ID
   */
  id: number,
  /**
   * 剧集标题
   */
  title: string,
  /**
   * 所属的课程ID
   */
  course_id: number,
  /**
   * 所属的课程
   */
  course: CourseResource
}

/**
 * 剧集列表元素数据
 */
export type EpisodesItemResource = {
  /**
   * 剧集ID
   */
  id: number,
  /**
   * 剧集标题
   */
  title: string,
}
