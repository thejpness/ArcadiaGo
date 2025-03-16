<script setup>
import { ref, computed } from "vue";
import { registerUser } from "@/api";
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import { usePasswordValidation } from "@/composables/usePasswordValidation"; // ✅ Import reusable password validation composable

const email = ref("");
const password = ref("");
const confirmPassword = ref("");
const error = ref("");
const success = ref("");
const router = useRouter();
const authStore = useAuthStore();

const showPassword = ref(false);
const showConfirmPassword = ref(false);

// ✅ Use password validation composable
const { passwordErrors, passwordMatchError, isPasswordValid } = usePasswordValidation(password, confirmPassword);

// ✅ Toggle password visibility
const passwordFieldType = computed(() => (showPassword.value ? "text" : "password"));
const confirmPasswordFieldType = computed(() => (showConfirmPassword.value ? "text" : "password"));

// ✅ Disable button unless all conditions are met
const isFormValid = computed(() => email.value && isPasswordValid.value);

// ✅ Register and Auto-Login
async function handleRegister() {
  if (!isFormValid.value) {
    error.value = "Please correct the errors before proceeding.";
    return;
  }

  try {
    await registerUser(email.value, password.value);
    success.value = "Account created! Logging you in...";

    // ✅ Auto-login after successful registration
    await authStore.login(email.value, password.value);

    // ✅ Redirect to Dashboard
    router.push("/dashboard");
  } catch (err) {
    error.value = err.message;
  }
}
</script>

<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100">
    <div class="w-full max-w-md p-8 bg-white rounded-lg shadow-lg">
      <h2 class="text-2xl font-semibold text-center">Register</h2>

      <div v-if="error" class="p-2 mt-2 text-red-600 bg-red-100 rounded">{{ error }}</div>
      <div v-if="success" class="p-2 mt-2 text-green-600 bg-green-100 rounded">{{ success }}</div>

      <input v-model="email" type="email" placeholder="Email" class="w-full p-2 mt-4 border rounded" />

      <!-- ✅ Password Field with Show/Hide -->
      <div class="relative">
        <input v-model="password" :type="passwordFieldType" placeholder="Password" class="w-full p-2 mt-2 border rounded pr-10" />
        <button type="button" @click="showPassword = !showPassword" class="absolute inset-y-0 right-2 flex items-center px-2 text-gray-600">
          {{ showPassword ? "👁️" : "🙈" }}
        </button>
      </div>

      <!-- ✅ Confirm Password Field with Show/Hide -->
      <div class="relative">
        <input v-model="confirmPassword" :type="confirmPasswordFieldType" placeholder="Confirm Password" class="w-full p-2 mt-2 border rounded pr-10" />
        <button type="button" @click="showConfirmPassword = !showConfirmPassword" class="absolute inset-y-0 right-2 flex items-center px-2 text-gray-600">
          {{ showConfirmPassword ? "👁️" : "🙈" }}
        </button>
      </div>

      <!-- ✅ Live Password Strength Validation -->
      <ul class="mt-2 text-sm text-gray-600">
        <li v-for="rule in passwordErrors" :key="rule" class="text-red-500">❌ {{ rule }}</li>
      </ul>

      <!-- ✅ Password Match Validation -->
      <p v-if="passwordMatchError" class="mt-2 text-sm text-red-500">{{ passwordMatchError }}</p>

      <!-- ✅ Button only enabled if all conditions are met -->
      <button 
        @click="handleRegister" 
        class="w-full p-2 mt-4 text-white bg-blue-500 rounded hover:bg-blue-600 disabled:opacity-50"
        :disabled="!isFormValid"
      >
        Register
      </button>

      <p class="mt-4 text-center">
        Already have an account?
        <router-link to="/login" class="text-blue-500">Login</router-link>
      </p>
    </div>
  </div>
</template>
