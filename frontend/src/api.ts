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
 * Fetches the authenticated user's details
 * @returns User data or null if not authenticated
 */
export async function fetchUser(): Promise<{ email: string } | null> {
  const token = localStorage.getItem("authToken");
  if (!token) {
    console.warn("⚠️ No auth token found. Skipping user fetch.");
    return null;
  }

  try {
    const response = await fetch(`${API_URL}/user`, {
      method: "GET",
      headers: { "Authorization": `Bearer ${token}` },
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
