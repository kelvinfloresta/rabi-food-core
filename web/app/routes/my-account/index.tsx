import { redirect, useLoaderData, Form } from "react-router"
import { useTranslation } from "react-i18next"
import type { Route } from "./+types/index"
import { resolveTenant } from "~/lib/tenant.server"
import { apiClient } from "~/lib/api.server"
import type { AccountProfile } from "~/components/subscription/types"
import { Button } from "~/components/ui/button"
import { PersonalInfoSection } from "./_components/personal-info-section"
import { AccessSection } from "./_components/access-section"
import { AddressSection } from "./_components/address-section"
import { BottomNav, BOTTOM_NAV_HEIGHT } from "~/components/bottom-nav"

export type AccountActionResult =
  | { intent: "save-personal-info"; ok: true }
  | { intent: "save-address"; ok: true }
  | { intent: string; error: string }

function parseAuthCookie(request: Request): string | null {
  const cookie = request.headers.get("Cookie") ?? ""
  return cookie.match(/auth_token=([^;]+)/)?.[1] ?? null
}

type SubscriptionListItem = Readonly<{ status: string }>
type SubscriptionPaginateResponse = Readonly<{ data: readonly SubscriptionListItem[] }>

export async function loader({ request }: Route.LoaderArgs) {
  const tenant = await resolveTenant(request)
  if (!tenant) throw redirect("/")

  const token = parseAuthCookie(request)
  if (!token) throw redirect("/")

  const api = apiClient(tenant.id)
  const [profileRes, subscriptionsRes] = await Promise.all([
    api.get<AccountProfile>("/user/me", undefined, token),
    api.get<SubscriptionPaginateResponse>("/subscription/", undefined, token),
  ])

  const profile: AccountProfile = profileRes.ok
    ? profileRes.data
    : { street: "", complement: "", neighborhood: "", city: "", state: "", zip: "", phone: "" }

  const hasActiveSubscription =
    subscriptionsRes.ok && subscriptionsRes.data.data.some((s) => s.status === "active")

  return { profile, hasActiveSubscription, token }
}

export async function action({ request }: Route.ActionArgs): Promise<AccountActionResult | Response> {
  const tenant = await resolveTenant(request)
  if (!tenant) return { intent: "unknown", error: "UNKNOWN" }

  const api = apiClient(tenant.id)
  const body = await request.json()
  const intent = body.intent as string

  if (intent === "save-personal-info") {
    const result = await api.patch("/user/me", { name: body.name, phone: body.phone }, body.token)
    if (!result.ok) return { intent: "save-personal-info", error: result.code }
    return { intent: "save-personal-info", ok: true }
  }

  if (intent === "save-address") {
    const result = await api.patch(
      "/user/me",
      {
        street: body.street,
        complement: body.complement,
        neighborhood: body.neighborhood,
        city: body.city,
        state: body.state,
        zip: body.zip,
      },
      body.token
    )
    if (!result.ok) return { intent: "save-address", error: result.code }
    return { intent: "save-address", ok: true }
  }

  return { intent, error: "UNKNOWN_INTENT" }
}

export default function MyAccount() {
  const { profile, hasActiveSubscription, token } = useLoaderData<typeof loader>()
  const { t } = useTranslation()

  return (
    <div className="min-h-screen bg-background text-foreground" style={{ paddingBottom: BOTTOM_NAV_HEIGHT + 16 }}>
      <header className="border-b px-6 py-4">
        <h1 className="text-xl font-semibold">{t("account.title")}</h1>
      </header>
      <main className="mx-auto max-w-2xl px-4 py-8 space-y-6">
        <PersonalInfoSection profile={profile} token={token} />
        <AccessSection email={profile.email} />
        <AddressSection
          profile={profile}
          token={token}
          hasActiveSubscription={hasActiveSubscription}
        />
        <Form method="post" action="/sign-out">
          <Button type="submit" variant="outline" className="w-full text-destructive border-destructive/40 hover:bg-destructive/5">
            {t("account.signOut")}
          </Button>
        </Form>
      </main>
      <BottomNav />
    </div>
  )
}
