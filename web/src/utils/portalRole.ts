export type PortalRole = 'parent' | 'mentor'

const ROLE_KEY = 'tm_portal_role'

export function getSavedPortalRole(): PortalRole | null {
  const v = localStorage.getItem(ROLE_KEY)
  if (v === 'parent' || v === 'mentor') return v
  return null
}

export function savePortalRole(role: PortalRole) {
  localStorage.setItem(ROLE_KEY, role)
}

export function loginPathForRole(role: PortalRole, redirect?: string) {
  const query: Record<string, string> = { role }
  if (redirect) query.redirect = redirect
  return { name: 'login', query }
}

export function homePathForAuth(isLoggedIn: boolean, isMentor: boolean, role?: PortalRole | null) {
  if (!isLoggedIn) return { name: 'welcome' }
  if (isMentor) return { name: 'mentor-home' }
  if (role === 'mentor') return { name: 'mentor-apply-status' }
  return { name: 'mentors' }
}
