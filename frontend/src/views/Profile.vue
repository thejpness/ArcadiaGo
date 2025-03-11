<script setup>
import { ref } from "vue";
import { useAuthStore } from "../stores/auth";
import { useRouter } from "vue-router";

const auth = useAuthStore();
const router = useRouter();

function logout() {
  auth.logout();
  router.push("/login");
}

// Mock user data (replace with actual backend API call later)
const user = ref({
  name: "John Doe",
  email: "john@example.com",
  joined: "2025-01-01"
});
</script>

<template>
  <div class="max-w-3xl mx-auto mt-10 p-6 bg-white rounded-lg shadow-lg">
    <h1 class="text-2xl font-bold">Profile</h1>
    <p class="text-gray-600">Manage your account details</p>

    <div class="mt-4 space-y-2">
      <p><strong>Name:</strong> {{ user.name }}</p>
      <p><strong>Email:</strong> {{ user.email }}</p>
      <p><strong>Joined:</strong> {{ user.joined }}</p>
    </div>

    <button @click="logout" class="w-full p-2 mt-6 text-white bg-red-500 rounded hover:bg-red-600">
      Logout
    </button>
  </div>
</template>
