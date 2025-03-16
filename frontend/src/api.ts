export const API_URL = "http://localhost:8080";

/**
 * Logs in the user and stores session via cookies
 * @param email - User email
 * @param password - User password
 * @returns Token string
 */
export async function loginUser(email: string, password: string): Promise<string> {
  const response = await fetch(`${API_URL}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include", // ✅ Ensures cookies are sent with the request
    body: JSON.stringify({ email, password }),
  });

  const text = await response.text();
  const data = text ? JSON.parse(text) : {};

  if (!response.ok) throw new Error(data.error || "Login failed");

  return data.token; // ✅ Ensure backend returns a token
}

/**
 * Registers a new user
 * @param email - New user's email
 * @param password - New user's password
 * @returns Success message
 */
export async function registerUser(email: string, password: string): Promise<string> {
  const response = await fetch(`${API_URL}/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include", // ✅ Ensures cookies are sent
    body: JSON.stringify({ email, password }),
  });

  const text = await response.text();
  const data = text ? JSON.parse(text) : {};

  if (!response.ok) throw new Error(data.error || "Registration failed");

  return data.message;
}

/**
 * Logs out the user by clearing session cookies
 * @returns Success message
 */
export async function logoutUser(): Promise<string> {
  const response = await fetch(`${API_URL}/logout`, {
    method: "POST",
    credentials: "include", // ✅ Ensures cookies are sent for logout
  });

  if (!response.ok) throw new Error("Logout failed");

  return "Logged out successfully";
}

/**
 * Fetches the authenticated user's details (now includes username)
 * @returns User data or null if not authenticated
 */
export async function fetchUser(): Promise<{ username: string; email: string; joined: string } | null> {
  try {
    const response = await fetch(`${API_URL}/user`, {
      method: "GET",
      credentials: "include",
    });

    if (!response.ok) {
      console.warn("⚠️ User is not authenticated:", response.status);
      return null; // ✅ Gracefully handle 401 errors instead of throwing
    }

    return await response.json();
  } catch (error) {
    console.error("❌ Failed to fetch user:", error);
    return null;
  }
}


export async function updateUsername(newUsername: string): Promise<void> {
  const response = await fetch(`${API_URL}/update-username`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    body: JSON.stringify({ new_username: newUsername }),
  });

  if (!response.ok) throw new Error("Failed to update username");
}

export async function updateEmail(newEmail: string): Promise<void> {
  const response = await fetch(`${API_URL}/update-email`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    body: JSON.stringify({ new_email: newEmail }),
  });

  if (!response.ok) throw new Error("Failed to request email change");
}

export async function updatePassword(oldPassword: string, newPassword: string): Promise<void> {
  const response = await fetch(`${API_URL}/update-password`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }),
  });

  if (!response.ok) throw new Error("Failed to update password");
}

export async function fetchActiveSessions(): Promise<{ id: string, userAgent: string, ipAddress: string }[]> {
  const response = await fetch(`${API_URL}/active-sessions`, {
    method: "GET",
    credentials: "include",
  });

  if (!response.ok) throw new Error("Failed to fetch active sessions");

  return response.json();
}

export async function logoutSession(sessionId: string): Promise<void> {
  const response = await fetch(`${API_URL}/logout-session`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    body: JSON.stringify({ session_id: sessionId }),
  });

  if (!response.ok) throw new Error("Failed to logout session");
}
