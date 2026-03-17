/**
 * 剧集数据
 */
export type EpisodesResource = {
  /** 剧集ID */
  id: number;
  /** 所属的剧集分组ID */
  group_id: number;
  /** 剧集标题 */
  title: string;
  /** 剧集文件路径 */
  file_path: string;
  /** 是否看完 */
  is_completed: boolean;
  /** 排序 */
  sort: number;
  /** 创建日期 */
  created_at: string;
};

/**
 * 剧集列表元素数据
 */
export type EpisodesItemResource = {
  /** 剧集ID */
  id: number;
  /** 剧集标题 */
  title: string;
};
