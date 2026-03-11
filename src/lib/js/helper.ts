/**
 * 根据传入的ID获取默认背景图片路径。
 *
 * @param id - 用于计算背景图片索引的数字ID。
 * @returns 返回对应索引的默认背景图片路径字符串。
 */
export function getDefaultBgImg(id: number): string {
  // 定义默认背景图片数组
  const defaultBgs = ["/img/bg-course-01.png", "/img/bg-course-02.png"];

  // 使用ID取模运算确定背景图片索引并返回对应路径
  return defaultBgs[id % defaultBgs.length] ?? "/img/bg-course-01.png";
}
