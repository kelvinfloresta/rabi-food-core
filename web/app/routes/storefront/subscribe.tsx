import { data } from "react-router"
import type { Route } from "./+types/subscribe"
import { resolveTenant } from "~/lib/tenant.server"
import { apiClient } from "~/lib/api.server"
import { createSessionHeaders, getAccessToken, type TokenPair } from "~/lib/session.server"
import type { WizardData } from "~/components/subscription/use-wizard"
import type { Product, SubscriptionConfig, UserProfile } from "~/components/subscription/types"

type OtpActionResult =
  | { intent: "send-otp"; ok: true }
  | { intent: "verify-otp"; ok: true }
  | { intent: "get-profile"; ok: true; profile: UserProfile }
  | { intent: "get-payment-profile"; ok: true; hasPaymentMethod: boolean; card?: { brand: string; last4: string; expMonth: number; expYear: number } }
  | { intent: "save-address"; ok: true }
  | { intent: string; error: string }

type SubscribeAndChargeResult =
  | { ok: true; subscriptionIds: readonly string[] }
  | { ok: true; requiresAction: true; clientSecret: string }
  | { error: string }

type SessionActionResult = ReturnType<typeof data<OtpActionResult>>

const UNAUTHENTICATED = "UNAUTHENTICATED"

export async function loader({ request, context }: Route.LoaderArgs) {
  const tenant = await resolveTenant(request)
  if (!tenant) throw new Response("Not found", { status: 404 })

  const api = apiClient(tenant.id)
  const token = getAccessToken(context)

  const [productsRes, configRes, profileRes] = await Promise.all([
    api.get<{ data: readonly Product[] }>("/product/", { isActive: "1" }),
    api.get<SubscriptionConfig>("/subscription/config"),
    token ? api.get<UserProfile>("/user/me", undefined, token) : Promise.resolve(null),
  ])

  if (!productsRes.ok) throw new Response("Not found", { status: 404 })
  if (!configRes.ok) throw new Response("Not found", { status: 404 })

  return {
    products: productsRes.data.data,
    config: configRes.data,
    isAuthenticated: token !== null,
    profile: profileRes?.ok ? profileRes.data : null,
  }
}

export async function action({ request, context }: Route.ActionArgs): Promise<OtpActionResult | SubscribeAndChargeResult | SessionActionResult> {
  const tenant = await resolveTenant(request)
  if (!tenant) return { intent: "unknown", error: "UNKNOWN" }

  const api = apiClient(tenant.id)
  const body = await request.json()
  const intent = body.intent as string | undefined

  if (intent === "send-otp") {
    const result = await api.post("/auth/otp", { email: body.email, name: body.name, phone: body.phone })
    if (!result.ok) return { intent: "send-otp", error: result.code }
    return { intent: "send-otp", ok: true }
  }

  if (intent === "verify-otp") {
    const result = await api.post<TokenPair>("/auth/otp/verify", {
      email: body.email,
      code: body.code,
    })
    if (!result.ok) return { intent: "verify-otp", error: result.code }
    return data<OtpActionResult>({ intent: "verify-otp", ok: true }, { headers: await createSessionHeaders(result.data) })
  }

  const token = getAccessToken(context)
  if (!token) return { intent: intent ?? "unknown", error: UNAUTHENTICATED }

  if (intent === "get-profile") {
    const result = await api.get<UserProfile>("/user/me", undefined, token)
    const profile = result.ok
      ? result.data
      : { street: "", complement: "", neighborhood: "", city: "", state: "", zip: "", phone: "" }
    return { intent: "get-profile", ok: true, profile }
  }

  if (intent === "get-payment-profile") {
    const result = await api.get<{ hasPaymentMethod: boolean; card?: { brand: string; last4: string; expMonth: number; expYear: number } }>("/payment/profile", undefined, token)
    return { intent: "get-payment-profile", ok: true, hasPaymentMethod: result.ok ? result.data.hasPaymentMethod : false, card: result.ok ? result.data.card : undefined }
  }

  if (intent === "save-address") {
    const result = await api.patch("/user/me", {
      phone: body.phone,
      street: body.street,
      complement: body.complement,
      neighborhood: body.neighborhood,
      city: body.city,
      state: body.state,
      zip: body.zip,
    }, token)
    if (!result.ok) return { intent: "save-address", error: result.code }
    return { intent: "save-address", ok: true }
  }

  // intent === "subscribe-and-charge"
  const checkout = body as WizardData & { checkoutKey: string; paymentToken?: string }

  const subResult = await api.post<readonly string[]>("/subscription/", {
    checkoutKey: checkout.checkoutKey,
    items: checkout.items,
    deliveryDays: checkout.deliveryDays,
    totalCycles: checkout.totalCycles,
  }, token)

  if (!subResult.ok) return { error: subResult.code }

  const paymentToken = checkout.paymentToken ?? ""

  for (const subscriptionId of subResult.data) {
    const chargeResult = await api.post<{ requiresAction?: boolean; clientSecret?: string } | null>(
      "/payment/charge",
      { subscriptionId, paymentToken },
      token,
    )
    if (!chargeResult.ok) return { error: chargeResult.code }
    if (chargeResult.data?.requiresAction && chargeResult.data?.clientSecret) {
      return { ok: true, requiresAction: true, clientSecret: chargeResult.data.clientSecret }
    }
  }

  return { ok: true, subscriptionIds: subResult.data }
}
