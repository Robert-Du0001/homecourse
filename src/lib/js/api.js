import { ElMessage } from 'element-plus';
import { useUserStore } from '@/store/user';

/**
 * 封装 Fetch 请求
 */
export async function request(method, url, data) {
  try {
    const userStore = useUserStore();

    const options = {
      method: method,
      headers: {
        "Content-Type": "application/json",
        ...(userStore.token ? { "Authorization": `Bearer ${userStore.token}` } : {})
      }
    };

    // 只有非 GET 请求才添加 body
    if (method.toUpperCase() !== 'GET' && data) {
      options.body = JSON.stringify(data);
    }

    const response = await fetch('/api'+url, options);
    
    // 解析 JSON
    const resData = await response.json();

    // 状态码非 200 处理
    if (response.status !== 200) {

      if (response.status === 401) {
        userStore.logout();
      }

      // 这里的 msg 对应你后端返回的 JSON 字段名
      ElMessage.error(resData.msg || '请求失败');
      return Promise.reject(resData); 
    }

    return resData; // 正常返回数据
  } catch (err) {
    console.error('Network Error:', err);
    ElMessage.error('网络连接异常');
    throw err; // 继续抛出，方便组件内处理 loading 状态
  }
}