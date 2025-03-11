export const API_URL = "http://localhost:8080";

/**
 * Logs in the user and stores session via cookies
 * @param email - User email
 * @param password - User password
 * @returns Success message
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

  return data.message;
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

  const text = await response.text();
  const data = text ? JSON.parse(text) : {};

  if (!response.ok) throw new Error(data.error || "Logout failed");

  return data.message;
}

/**
 * Fetches the authenticated user's details
 * @returns User data (Modify this to match backend response)
 */
export async function fetchUser(): Promise<{ email: string }> {
  const response = await fetch(`${API_URL}/user`, {
    method: "GET",
    credentials: "include", // ✅ Ensures cookies are included in request
  });

  const text = await response.text();
  const data = text ? JSON.parse(text) : {};

  if (!response.ok) throw new Error(data.error || "Failed to fetch user data");

  return data;
}
