<script setup>
import { useAuthStore } from "../stores/auth";
import { useRouter } from "vue-router";

const auth = useAuthStore();
const router = useRouter();

async function logout() {
  try {
    await auth.logout();
    router.push("/login"); // ✅ Redirect to login after logout
  } catch (error) {
    console.error("Logout failed:", error);
  }
}
</script>

<template>
  <div class="flex h-screen bg-gray-100">
    <!-- Sidebar -->
    <aside class="w-64 bg-white shadow-lg">
      <div class="p-6 text-lg font-semibold border-b">ArcadiaGo</div>
      <nav class="p-4">
        <ul>
          <li>
            <router-link to="/dashboard" class="flex items-center p-3 rounded hover:bg-gray-200">
              🏠 Dashboard
            </router-link>
          </li>
          <li>
            <router-link to="/profile" class="flex items-center p-3 rounded hover:bg-gray-200">
              👤 Profile
            </router-link>
          </li>
          <li>
            <router-link to="/settings" class="flex items-center p-3 rounded hover:bg-gray-200">
              ⚙️ Settings
            </router-link>
          </li>
        </ul>
      </nav>
      <button @click="logout" class="w-full p-3 mt-4 text-white bg-red-500 rounded hover:bg-red-600">
        Logout
      </button>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 p-6">
      <router-view></router-view>
    </main>
  </div>
</template>
