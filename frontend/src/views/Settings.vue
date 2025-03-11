<script setup>
import { ref } from "vue";

const settings = ref({
  notifications: true,
  darkMode: false
});

function toggleSetting(setting) {
  settings.value[setting] = !settings.value[setting];
}
</script>

<template>
  <div class="max-w-3xl mx-auto mt-10 p-6 bg-white rounded-lg shadow-lg">
    <h1 class="text-2xl font-bold">Settings</h1>
    <p class="text-gray-600">Customize your preferences</p>

    <div class="mt-4 space-y-4">
      <label class="flex items-center space-x-3 cursor-pointer">
        <input type="checkbox" v-model="settings.notifications" class="w-5 h-5">
        <span>Enable Notifications</span>
      </label>

      <label class="flex items-center space-x-3 cursor-pointer">
        <input type="checkbox" v-model="settings.darkMode" class="w-5 h-5">
        <span>Enable Dark Mode</span>
      </label>
    </div>
  </div>
</template>
