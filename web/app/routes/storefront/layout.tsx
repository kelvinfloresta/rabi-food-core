import { Outlet, data } from "react-router"
import type { Route } from "./+types/layout"
import { getTenantSlug, type Tenant } from "~/lib/tenant.server"
import { localeCookie } from "~/i18n.server"
import { apiClient } from "~/lib/api.server"

export async function loader({ request }: Route.LoaderArgs) {
  const slug = getTenantSlug(request)

  // No slug means dashboard subdomain — skip storefront, let dashboard layout handle it
  if (!slug) return data({ tenant: null })

  const result = await apiClient().get<Tenant>(`/tenant/slug/${slug}`)
  if (!result.ok) throw new Response(result.code, { status: 500 })

  const headers = new Headers()

  const hasLangCookie = await localeCookie.parse(request.headers.get("Cookie"))
  if (!hasLangCookie) {
    headers.set("Set-Cookie", await localeCookie.serialize(result.data.language))
  }

  return data({ tenant: result.data }, { headers })
}

export default function StorefrontLayout() {
  return <Outlet />
}
