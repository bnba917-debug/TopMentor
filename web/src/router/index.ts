import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/welcome' },
    {
      path: '/welcome',
      name: 'welcome',
      component: () => import('@/views/WelcomeView.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { guestOnly: true },
    },
    {
      path: '/mentors',
      name: 'mentors',
      component: () => import('@/views/MentorsView.vue'),
    },
    {
      path: '/mentors/:id',
      name: 'mentor-detail',
      component: () => import('@/views/MentorDetailView.vue'),
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('@/views/ProfileView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/packages',
      name: 'packages',
      component: () => import('@/views/PackagesView.vue'),
    },
    {
      path: '/orders',
      name: 'orders',
      component: () => import('@/views/OrdersView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/reports/:orderId',
      name: 'report',
      component: () => import('@/views/ReportView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/classroom/:orderId',
      name: 'classroom',
      component: () => import('@/views/ClassroomView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/mentor',
      name: 'mentor-home',
      component: () => import('@/views/mentor/MentorHomeView.vue'),
      meta: { requiresAuth: true, requiresMentor: true },
    },
    {
      path: '/mentor/orders',
      name: 'mentor-orders',
      component: () => import('@/views/mentor/MentorOrdersView.vue'),
      meta: { requiresAuth: true, requiresMentor: true },
    },
    {
      path: '/mentor/orders/:orderId/report',
      name: 'mentor-report',
      component: () => import('@/views/mentor/MentorReportView.vue'),
      meta: { requiresAuth: true, requiresMentor: true },
    },
    {
      path: '/mentor/slots',
      name: 'mentor-slots',
      component: () => import('@/views/mentor/MentorSlotsView.vue'),
      meta: { requiresAuth: true, requiresMentor: true },
    },
    {
      path: '/mentor/wallet',
      name: 'mentor-wallet',
      component: () => import('@/views/mentor/MentorWalletView.vue'),
      meta: { requiresAuth: true, requiresMentor: true },
    },
    {
      path: '/mentor/profile',
      name: 'mentor-profile',
      component: () => import('@/views/mentor/MentorProfileView.vue'),
      meta: { requiresAuth: true, requiresMentor: true },
    },
    {
      path: '/mentor/apply',
      name: 'mentor-apply',
      component: () => import('@/views/mentor/MentorApplyView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/mentor/apply/status',
      name: 'mentor-apply-status',
      component: () => import('@/views/mentor/MentorApplyStatusView.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.isLoggedIn) {
    const role = to.fullPath.startsWith('/mentor') ? 'mentor' : 'parent'
    return { name: 'login', query: { redirect: to.fullPath, role } }
  }
  if (to.meta.requiresMentor && !auth.isMentor) {
    return { name: 'mentor-apply-status' }
  }
  if (to.meta.guestOnly && auth.isLoggedIn) {
    return auth.isMentor ? { name: 'mentor-home' } : { name: 'mentors' }
  }
})

export default router
