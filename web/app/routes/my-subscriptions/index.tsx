import { redirect, useLoaderData } from "react-router"
import { useTranslation } from "react-i18next"
import { Package } from "lucide-react"
import type { Route } from "./+types/index"
import { resolveTenant } from "~/lib/tenant.server"
import { apiClient } from "~/lib/api.server"
import { getAccessToken } from "~/lib/session.server"
import { Button } from "~/components/ui/button"
import { BottomNav, BOTTOM_NAV_HEIGHT } from "~/components/bottom-nav"

type SubscriptionStatus = "active" | "cancelled" | string

type SubscriptionItem = Readonly<{
  id: string
  status: SubscriptionStatus
  totalCycles: number
  remainingCycles: number
  paymentAmount: number
  autoRenew: boolean
  createdAt: string
}>

type SubscriptionsResponse = Readonly<{
  data: readonly SubscriptionItem[]
  maxPages: number
}>

function formatPrice(cents: number): string {
  return (cents / 100).toLocaleString("pt-BR", { style: "currency", currency: "BRL" })
}

export async function loader({ request, context }: Route.LoaderArgs) {
  const tenant = await resolveTenant(request)
  if (!tenant) throw redirect("/")

  const token = getAccessToken(context)
  if (!token) throw redirect("/")

  const api = apiClient(tenant.id)
  const res = await api.get<SubscriptionsResponse>("/subscription/", undefined, token)

  const items: readonly SubscriptionItem[] = res.ok ? res.data.data : []

  const active = items.filter((s) => s.status === "active")
  const inactive = items.filter((s) => s.status !== "active")

  return { active, inactive }
}

type SubscriptionCardProps = Readonly<{
  item: SubscriptionItem
}>

function SubscriptionCard({ item }: SubscriptionCardProps) {
  const { t } = useTranslation()

  const statusLabel =
    item.status === "active"
      ? t("subscriptions.status.active")
      : item.status === "cancelled"
        ? t("subscriptions.status.cancelled")
        : t("subscriptions.status.unknown")

  const statusClass =
    item.status === "active"
      ? "text-green-700 bg-green-50 border-green-200"
      : "text-muted-foreground bg-muted border-border"

  return (
    <div className="rounded-xl border bg-card p-4 space-y-3">
      <div className="flex items-start justify-between gap-2">
        <div className="flex items-center gap-2 min-w-0">
          <Package className="size-4 shrink-0 text-muted-foreground" />
          <span className="text-sm font-medium truncate">{item.id}</span>
        </div>
        <span className={`shrink-0 rounded-full border px-2 py-0.5 text-xs font-medium ${statusClass}`}>
          {statusLabel}
        </span>
      </div>
      <dl className="grid grid-cols-2 gap-x-4 gap-y-2 text-sm">
        <div>
          <dt className="text-xs text-muted-foreground">{t("subscriptions.remainingCycles", { remaining: item.remainingCycles, total: item.totalCycles })}</dt>
        </div>
        <div className="text-right">
          <dd className="font-semibold text-primary">{formatPrice(item.paymentAmount)}</dd>
          <dt className="text-xs text-muted-foreground">{t("subscriptions.perDelivery")}</dt>
        </div>
      </dl>
      {item.status === "active" && item.autoRenew && (
        <Button variant="outline" size="sm" disabled className="w-full">
          {t("subscriptions.cancelRenewal")}
        </Button>
      )}
    </div>
  )
}

export default function MySubscriptions() {
  const { active, inactive } = useLoaderData<typeof loader>()
  const { t } = useTranslation()

  const hasAny = active.length > 0 || inactive.length > 0

  return (
    <div className="min-h-screen bg-background text-foreground" style={{ paddingBottom: BOTTOM_NAV_HEIGHT + 16 }}>
      <header className="border-b px-6 py-4">
        <h1 className="text-xl font-semibold">{t("subscriptions.title")}</h1>
      </header>

      {!hasAny && (
        <div className="flex flex-col items-center justify-center gap-3 px-6 py-20 text-center">
          <Package className="size-10 text-muted-foreground" />
          <p className="text-sm text-muted-foreground">{t("subscriptions.empty")}</p>
        </div>
      )}

      {hasAny && (
        <main className="mx-auto max-w-2xl px-4 py-6 space-y-8">
          {active.length > 0 && (
            <section className="space-y-3">
              <h2 className="text-sm font-medium text-muted-foreground uppercase tracking-wide">
                {t("subscriptions.status.active")}
              </h2>
              {active.map((item) => (
                <SubscriptionCard key={item.id} item={item} />
              ))}
            </section>
          )}
          {inactive.length > 0 && (
            <section className="space-y-3">
              <h2 className="text-sm font-medium text-muted-foreground uppercase tracking-wide">
                {t("subscriptions.status.cancelled")}
              </h2>
              {inactive.map((item) => (
                <SubscriptionCard key={item.id} item={item} />
              ))}
            </section>
          )}
        </main>
      )}

      <BottomNav />
    </div>
  )
}
