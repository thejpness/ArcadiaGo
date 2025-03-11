<script setup>
import { ref } from "vue";
import { useAuthStore } from "../stores/auth";
import { useRouter } from "vue-router";

const email = ref("");
const password = ref("");
const error = ref("");
const auth = useAuthStore();
const router = useRouter();

async function handleLogin() {
  try {
    await auth.login(email.value, password.value);
    router.push("/dashboard");
  } catch (err) {
    error.value = err.message;
  }
}
</script>

<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100">
    <div class="w-full max-w-md p-8 bg-white rounded-lg shadow-lg">
      <h2 class="text-2xl font-semibold text-center">Login</h2>

      <div v-if="error" class="p-2 mt-2 text-red-600 bg-red-100 rounded">{{ error }}</div>

      <input v-model="email" type="email" placeholder="Email" class="w-full p-2 mt-4 border rounded" />
      <input v-model="password" type="password" placeholder="Password" class="w-full p-2 mt-2 border rounded" />
      <button @click="handleLogin" class="w-full p-2 mt-4 text-white bg-blue-500 rounded hover:bg-blue-600">
        Login
      </button>

      <p class="mt-4 text-center">
        Don't have an account?
        <router-link to="/register" class="text-blue-500">Register</router-link>
      </p>
    </div>
  </div>
</template>
