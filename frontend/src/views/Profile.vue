<script setup>
import { ref, computed, onMounted } from "vue";
import { useAuthStore } from "@/stores/auth";
import { useRouter } from "vue-router";
import { fetchUser, updateEmail, updatePassword, updateUsername, fetchActiveSessions, logoutSession } from "@/api";
import { usePasswordValidation } from "@/composables/usePasswordValidation"; // ✅ Import reusable password validation composable

const auth = useAuthStore();
const router = useRouter();

const user = ref({
  username: "",
  email: "",
  joined: "",
});

const sessions = ref([]); // Active sessions

const newUsername = ref("");
const newEmail = ref("");
const oldPassword = ref("");
const newPassword = ref("");
const confirmPassword = ref("");
const successMessage = ref("");
const errorMessage = ref("");

const showOldPassword = ref(false);
const showNewPassword = ref(false);
const showConfirmPassword = ref(false);

// ✅ Use password validation composable
const { passwordErrors, passwordMatchError, isPasswordValid } = usePasswordValidation(newPassword, confirmPassword);

// ✅ Fetch user data on component mount
onMounted(async () => {
  try {
    const userData = await fetchUser();
    if (userData) {
      user.value = userData;
      const sessionData = await fetchActiveSessions();
      sessions.value = sessionData || [];
    }
  } catch (error) {
    console.error("Failed to load profile:", error);
  }
});

// ✅ Update Username
async function handleUpdateUsername() {
  try {
    await updateUsername(newUsername.value);
    successMessage.value = "Username updated successfully";
    user.value.username = newUsername.value;
    newUsername.value = "";
  } catch (error) {
    errorMessage.value = error.message;
  }
}

// ✅ Update Email
async function handleUpdateEmail() {
  try {
    await updateEmail(newEmail.value);
    successMessage.value = "Email change request sent! Check your inbox.";
    newEmail.value = "";
  } catch (error) {
    errorMessage.value = error.message;
  }
}

// ✅ Update Password
async function handleUpdatePassword() {
  if (!isPasswordValid.value) {
    errorMessage.value = "Please correct password errors before updating.";
    return;
  }
  try {
    await updatePassword(oldPassword.value, newPassword.value);
    successMessage.value = "Password updated successfully";
    oldPassword.value = newPassword.value = confirmPassword.value = "";
  } catch (error) {
    errorMessage.value = error.message;
  }
}

// ✅ Logout from a specific session
async function handleLogoutSession(sessionId) {
  try {
    await logoutSession(sessionId);
    sessions.value = sessions.value.filter(session => session.id !== sessionId);
    successMessage.value = "Logged out of session.";
  } catch (error) {
    errorMessage.value = error.message;
  }
}

// ✅ Logout completely
function logout() {
  auth.logout();
  router.push("/login");
}
</script>

<template>
  <div class="max-w-3xl mx-auto mt-10 p-6 bg-white rounded-lg shadow-lg">
    <h1 class="text-2xl font-bold">Profile</h1>
    <p class="text-gray-600">Manage your account details</p>

    <div v-if="successMessage" class="mt-2 p-2 bg-green-100 text-green-700 rounded">{{ successMessage }}</div>
    <div v-if="errorMessage" class="mt-2 p-2 bg-red-100 text-red-700 rounded">{{ errorMessage }}</div>

    <!-- User Info -->
    <div class="mt-4 space-y-2">
      <p><strong>Username:</strong> {{ user.username }}</p>
      <p><strong>Email:</strong> {{ user.email }}</p>
      <p><strong>Joined:</strong> {{ user.joined }}</p>
    </div>

    <!-- Update Username -->
    <div class="mt-4">
      <input v-model="newUsername" placeholder="New Username" class="border p-2 w-full rounded" />
      <button @click="handleUpdateUsername" class="mt-2 w-full p-2 text-white bg-blue-500 rounded hover:bg-blue-600">
        Update Username
      </button>
    </div>

    <!-- Update Email -->
    <div class="mt-4">
      <input v-model="newEmail" placeholder="New Email" class="border p-2 w-full rounded" />
      <button @click="handleUpdateEmail" class="mt-2 w-full p-2 text-white bg-yellow-500 rounded hover:bg-yellow-600">
        Change Email
      </button>
    </div>

    <!-- Update Password -->
    <div class="mt-4">
      <h2 class="text-xl font-bold">Update Password</h2>

      <!-- Old Password -->
      <div class="relative">
        <input v-model="oldPassword" :type="showOldPassword ? 'text' : 'password'" placeholder="Old Password" class="border p-2 w-full rounded pr-10" />
        <button type="button" @click="showOldPassword = !showOldPassword" class="absolute inset-y-0 right-2 flex items-center px-2 text-gray-600">
          {{ showOldPassword ? "👁️" : "🙈" }}
        </button>
      </div>

      <!-- New Password -->
      <div class="relative mt-2">
        <input v-model="newPassword" :type="showNewPassword ? 'text' : 'password'" placeholder="New Password" class="border p-2 w-full rounded pr-10" />
        <button type="button" @click="showNewPassword = !showNewPassword" class="absolute inset-y-0 right-2 flex items-center px-2 text-gray-600">
          {{ showNewPassword ? "👁️" : "🙈" }}
        </button>
      </div>

      <!-- Confirm Password -->
      <div class="relative mt-2">
        <input v-model="confirmPassword" :type="showConfirmPassword ? 'text' : 'password'" placeholder="Confirm Password" class="border p-2 w-full rounded pr-10" />
        <button type="button" @click="showConfirmPassword = !showConfirmPassword" class="absolute inset-y-0 right-2 flex items-center px-2 text-gray-600">
          {{ showConfirmPassword ? "👁️" : "🙈" }}
        </button>
      </div>

      <!-- Password Strength Validation -->
      <ul class="mt-2 text-sm text-gray-600">
        <li v-for="rule in passwordErrors" :key="rule" class="text-red-500">❌ {{ rule }}</li>
      </ul>

      <!-- Password Match Validation -->
      <p v-if="passwordMatchError" class="mt-2 text-sm text-red-500">{{ passwordMatchError }}</p>

      <!-- Update Password Button -->
      <button 
        @click="handleUpdatePassword" 
        class="mt-2 w-full p-2 text-white bg-green-500 rounded hover:bg-green-600 disabled:opacity-50"
        :disabled="!isPasswordValid"
      >
        Update Password
      </button>
    </div>

    <!-- Active Sessions -->
    <div class="mt-4">
      <h2 class="text-xl font-bold">Active Sessions</h2>
      <ul v-if="sessions.length">
        <li v-for="session in sessions" :key="session.id" class="flex justify-between items-center p-2 bg-gray-100 rounded mt-2">
          <span>{{ session.userAgent }} ({{ session.ipAddress }})</span>
          <button @click="handleLogoutSession(session.id)" class="text-red-500">Logout</button>
        </li>
      </ul>
      <p v-else class="text-gray-600">No active sessions.</p>
    </div>

    <!-- Logout -->
    <button @click="logout" class="w-full p-2 mt-6 text-white bg-red-500 rounded hover:bg-red-600">
      Logout
    </button>
  </div>
</template>
