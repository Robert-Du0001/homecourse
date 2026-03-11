import { defineStore } from "pinia";

import type { UserResource } from "@/types/user";

import router from "@/router";
import { UserRole } from "@/types/user";

export const useUserStore = defineStore("user", {
  state: (): UserResource => ({
    id: 0,
    token: localStorage.getItem("token") || "",
    name: "",
    role: UserRole.GUEST,
  }),
  actions: {
    setToken(newToken: string) {
      this.token = newToken;
      localStorage.setItem("token", newToken);
    },
    logout() {
      localStorage.removeItem("token");
      this.$reset();
      router.replace({ name: "Login" });
    },
  },
});
