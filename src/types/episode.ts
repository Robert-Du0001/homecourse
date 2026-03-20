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
  // /** 总时长（秒） */
  // duration: number;
  // /** 已观看时长（秒） */
  // watched_duration: number;
  /** 所属剧集分组 */
  group: {
    /** 剧集分组ID */
    id: number;
    /** 剧集分组标题 */
    name: string;
    /** 所属课程ID */
    course_id: number;
    /** 所属课程 */
    course: {
      /** 课程ID */
      id: number;
      /** 课程标题 */
      title: string;
    };
  };
  attachments: {
    /** 附件ID */
    id: number;
    /** 附件所属剧集ID */
    episode_id: number;
    /** 附件类型 */
    name: string;
  }[];
};

/**
 * 剧集列表元素数据
 */
export type EpisodesItemResource = {
  /** 剧集ID */
  id: number;
  /** 所属的剧集分组ID */
  group_id: number;
  /** 剧集标题 */
  title: string;
  /** 剧集文件路径 */
  file_path: string;
  /** 总时长（秒） */
  duration: number;
  /** 已观看时长（秒） */
  watched_duration: number;
  /** 排序 */
  sort: number;
  /** 创建日期 */
  created_at: string;
};
