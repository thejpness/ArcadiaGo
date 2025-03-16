import { ref, computed, type Ref } from "vue";

/**
 * Password validation composable for ensuring password security and confirmation match.
 * 
 * @param passwordRef - A Vue ref containing the user's password.
 * @param confirmPasswordRef - A Vue ref containing the confirmed password.
 * @returns Computed properties for password errors, match validation, and overall validity.
 */
export function usePasswordValidation(passwordRef: Ref<string>, confirmPasswordRef: Ref<string>) {
  // ✅ Password security rules
  const passwordRules = {
    minLength: 8,
    hasUppercase: /[A-Z]/,
    hasLowercase: /[a-z]/,
    hasNumber: /\d/,
    hasSpecial: /[@$!%*?&]/,
  };

  // ✅ Live validation for password security
  const passwordErrors = computed(() => {
    const errors: string[] = [];
    if (passwordRef.value.length < passwordRules.minLength) errors.push("At least 8 characters.");
    if (!passwordRules.hasUppercase.test(passwordRef.value)) errors.push("At least 1 uppercase letter.");
    if (!passwordRules.hasLowercase.test(passwordRef.value)) errors.push("At least 1 lowercase letter.");
    if (!passwordRules.hasNumber.test(passwordRef.value)) errors.push("At least 1 number.");
    if (!passwordRules.hasSpecial.test(passwordRef.value)) errors.push("At least 1 special character (@$!%*?&).");
    return errors;
  });

  // ✅ Live validation for password match
  const passwordMatchError = computed(() => {
    return passwordRef.value !== confirmPasswordRef.value ? "Passwords do not match." : "";
  });

  // ✅ Computed property to check if the password meets all criteria
  const isPasswordValid = computed(() => {
    return (
      passwordRef.value.length >= passwordRules.minLength &&
      passwordRules.hasUppercase.test(passwordRef.value) &&
      passwordRules.hasLowercase.test(passwordRef.value) &&
      passwordRules.hasNumber.test(passwordRef.value) &&
      passwordRules.hasSpecial.test(passwordRef.value) &&
      passwordRef.value === confirmPasswordRef.value
    );
  });

  return { passwordErrors, passwordMatchError, isPasswordValid };
}
