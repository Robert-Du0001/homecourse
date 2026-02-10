/**
 * 用户角色枚举
 */
export enum UserRole {
  /** 普通用户 */
  GUEST = 0,
  /** 管理员 */
  ADMIN = 1,
}

/**
 * 用户数据
 */
export type UserResource = {
  /**
   * 用户名
   */
  name: string,
  /**
   * 角色
   */
  role: UserRole,
  /**
   * JWT登录凭证
   */
  token: string,
};
