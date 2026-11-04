import { apiClient } from "~/lib/api.server"
import type { Result } from "./result"

const DASHBOARD_SLUGS = new Set(["app", "www"])

export type Tenant = {
  id: string
  name: string
  slug: string
  language: string
}

export function getTenantSlug(request: Request): string | null {
  const host = request.headers.get("host")
  if (!host) return null
  const indexOfDot = host.indexOf(".")
  
  if (indexOfDot === -1) return null
  const slug = host.slice(0, indexOfDot)

  // bare localhost or app.* subdomain → dashboard
  if (slug === "localhost" || DASHBOARD_SLUGS.has(slug)) {
    return null
  }

  return slug
}

export function isDashboard(request: Request): boolean {
  return getTenantSlug(request) === null
}

export async function resolveTenant(request: Request): Promise<Tenant | null> {
  const slug = getTenantSlug(request)
  if (!slug) return null
  const result = await apiClient().get<Tenant>(`/tenant/slug/${slug}`)
  return result.ok ? result.data : null
}
