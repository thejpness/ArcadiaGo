import { defineStore } from "pinia";
import { ref, computed, watchEffect } from "vue";
import { loginUser, logoutUser, fetchUser } from "../api";
import { useRouter } from "vue-router";

export const useAuthStore = defineStore("auth", () => {
  const isAuthenticated = ref<boolean>(false);
  const userEmail = ref<string | null>(null);
  const loading = ref<boolean>(false);
  const errorMessage = ref<string | null>(null);
  const router = useRouter();

  // ✅ Computed values for cleaner template binding
  const isLoggedIn = computed(() => isAuthenticated.value);
  const hasError = computed(() => errorMessage.value !== null);

  /**
   * Logs in the user, updates authentication state, and redirects.
   */
  async function login(email: string, password: string) {
    loading.value = true;
    errorMessage.value = null;

    try {
      await loginUser(email, password);
      await loadUser(); // ✅ Fetch user data after login
      router.push("/dashboard"); // ✅ Redirect to dashboard on success
    } catch (error) {
      console.error("❌ Login error:", error);
      errorMessage.value = error instanceof Error ? error.message : "Login failed";
      throw error;
    } finally {
      loading.value = false;
    }
  }

  /**
   * Logs out the user, clears authentication state, and redirects.
   */
  async function logout() {
    loading.value = true;
    errorMessage.value = null;

    try {
      await logoutUser();
      isAuthenticated.value = false;
      userEmail.value = null;
      router.push("/login"); // ✅ Redirect to login on logout
    } catch (error) {
      console.error("❌ Logout error:", error);
      errorMessage.value = "Logout failed. Try again.";
    } finally {
      loading.value = false;
    }
  }

  /**
   * Fetches the authenticated user's data and updates the store.
   * Called when the app loads to persist session.
   */
  async function loadUser() {
    loading.value = true;
    errorMessage.value = null;

    try {
      const user = await fetchUser();
      userEmail.value = user.email;
      isAuthenticated.value = true;
    } catch (error) {
      console.warn("⚠️ Not authenticated:", error);
      isAuthenticated.value = false;
      userEmail.value = null;
    } finally {
      loading.value = false;
    }
  }

  // ✅ Auto-fetch user data when the store is initialized
  watchEffect(() => {
    loadUser();
  });

  return {
    isAuthenticated,
    isLoggedIn,
    userEmail,
    loading,
    errorMessage,
    hasError,
    login,
    logout,
    loadUser,
  };
});
