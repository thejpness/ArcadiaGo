import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/stores/auth"; // ✅ Use alias for better maintainability

// ✅ Lazy-load views for better performance (Code-splitting)
const Login = () => import("@/views/Login.vue");
const Register = () => import("@/views/Register.vue");
const Dashboard = () => import("@/views/Dashboard.vue");
const Profile = () => import("@/views/Profile.vue");
const Settings = () => import("@/views/Settings.vue");

const routes = [
  { path: "/", redirect: "/dashboard" }, // ✅ Redirect root to Dashboard (if logged in)
  { path: "/login", component: Login, meta: { requiresAuth: false, guestOnly: true } },
  { path: "/register", component: Register, meta: { requiresAuth: false, guestOnly: true } },
  { path: "/dashboard", component: Dashboard, meta: { requiresAuth: true } },
  { path: "/profile", component: Profile, meta: { requiresAuth: true } },
  { path: "/settings", component: Settings, meta: { requiresAuth: true } },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});

// ✅ Navigation Guard to Protect Routes
router.beforeEach(async (to, from, next) => {
  const auth = useAuthStore(); // ✅ Access authentication store
  await auth.loadUser(); // ✅ Ensure authentication status is updated before checking

  // 🔒 Redirect unauthorized users to login
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return next("/login");
  }

  // 🚀 Prevent authenticated users from accessing login/register
  if (to.meta.guestOnly && auth.isAuthenticated) {
    return next("/dashboard");
  }

  next(); // ✅ Allow navigation
});
