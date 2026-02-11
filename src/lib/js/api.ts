import { ElLoading } from 'element-plus';
import { useUserStore } from '@/store/user';

/**
 * 发送请求允许的方法
 */
export type Method = 'GET' | 'POST' | 'PUT' | 'DELETE';

/**
 * API响应结果
 */
export type ApiResponse<T = null> = {
  /**
   * 响应消息
   */
  msg: string;
  /**
   * 响应内容
   */
  data: T;
};

/**
 * 响应返回的报错数据
 */
export type CatchData = {
  /**
   * 状态码
   */
  status: number;
  /**
   * 错误消息
   */
  msg: string;
  /**
   * 附带的数据
   */
  data?: object;
};

/**
 * 封装 Fetch 请求
 */
export async function request<T = null>(
  method: Method | 'DELETE',
  url: string,
  data?: object,
) {
  return new Promise<ApiResponse<T> & { status: string }>((resolve, reject) => {
    const loading = ElLoading.service({
      lock: true,
      text: '加载中',
      background: 'rgba(0, 0, 0, 0.7)',
    });

    const userStore = useUserStore();

    const options: RequestInit = {
      method,
      headers: {
        'Content-Type': 'application/json',
        ...(userStore.token
          ? { Authorization: `Bearer ${userStore.token}` }
          : {}),
      },
    };

    // 只有非 GET 请求才添加 body
    if (method !== 'GET' && data) {
      options.body = JSON.stringify(data);
    }

    fetch('/api' + url, options)
      .then((res) => {
        if (res.status === 401) {
          userStore.logout();
          reject({
            status: res.status,
            msg: '无请求权限！',
            data,
          });
        } else {
          // 这里需要在后续使用 res 变量，所以就直接在里面解析json了
          return res
            .json()
            .then((resData: ApiResponse<T>) => {
              const { msg, data } = resData;

              if (res.status === 200) {
                resolve({ status: 'ok', msg, data });
              } else {
                reject({
                  status: res.status,
                  msg: msg || '发生错误！',
                  data,
                });
              }
            })
            .catch((err) => {
              console.error('解析json数据失败:', err);
              reject({ msg: '解析json数据失败' });
            });
        }
      })
      .catch((err) => {
        console.error('发起请求失败:', err);
        reject({ msg: '发起请求失败' });
      })
      .finally(() => {
        loading.close();
      });
  });
}
