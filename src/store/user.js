import { defineStore } from 'pinia';
import router from '@/router';

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    name: '', 
  }),
  actions: {
    setToken(newToken) {
      this.token = newToken;
      localStorage.setItem('token', newToken);
    },
    logout() {
      localStorage.removeItem('token');
      this.$reset();
      
      router.replace({name: 'Login'});
    }
  },
});