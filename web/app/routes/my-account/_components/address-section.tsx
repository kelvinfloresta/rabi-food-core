import { useState } from "react"
import { useFetcher } from "react-router"
import { useTranslation } from "react-i18next"
import { AlertTriangle } from "lucide-react"
import { Button } from "~/components/ui/button"
import { Input } from "~/components/ui/input"
import type { AccountProfile } from "~/components/subscription/types"
import type { AccountActionResult } from "../index"

type AddressSectionProps = Readonly<{
  profile: AccountProfile
  token: string
  hasActiveSubscription: boolean
}>

type AddressFields = {
  street: string
  complement: string
  neighborhood: string
  city: string
  state: string
  zip: string
}

type AddressErrors = Readonly<Partial<Record<keyof AddressFields, string>>>
type MutableAddressErrors = Partial<Record<keyof AddressFields, string>>

const ZIP_PLACEHOLDER = "00000-000"

export function AddressSection({ profile, token, hasActiveSubscription }: AddressSectionProps) {
  const { t } = useTranslation()
  const fetcher = useFetcher<AccountActionResult>()
  const [editing, setEditing] = useState(false)
  const [fields, setFields] = useState<AddressFields>({
    street: profile.street ?? "",
    complement: profile.complement ?? "",
    neighborhood: profile.neighborhood ?? "",
    city: profile.city ?? "",
    state: profile.state ?? "",
    zip: profile.zip ?? "",
  })
  const [errors, setErrors] = useState<AddressErrors>({})

  const isSaving = fetcher.state !== "idle"

  function handleEdit() {
    setFields({
      street: profile.street ?? "",
      complement: profile.complement ?? "",
      neighborhood: profile.neighborhood ?? "",
      city: profile.city ?? "",
      state: profile.state ?? "",
      zip: profile.zip ?? "",
    })
    setErrors({})
    setEditing(true)
  }

  function handleCancel() {
    setEditing(false)
    setErrors({})
  }

  function set(field: keyof AddressFields) {
    return (e: React.ChangeEvent<HTMLInputElement>) =>
      setFields((prev) => ({ ...prev, [field]: e.target.value }))
  }

  function validate(): boolean {
    const next: MutableAddressErrors = {}
    if (!fields.street.trim()) next.street = t("account.sections.address.streetRequired")
    if (!fields.neighborhood.trim()) next.neighborhood = t("account.sections.address.neighborhoodRequired")
    if (!fields.city.trim()) next.city = t("account.sections.address.cityRequired")
    if (!fields.state.trim()) next.state = t("account.sections.address.stateRequired")
    if (!fields.zip.trim()) next.zip = t("account.sections.address.zipRequired")
    setErrors(next)
    return Object.keys(next).length === 0
  }

  function handleSave() {
    if (!validate()) return
    fetcher.submit(
      { intent: "save-address", ...fields, token },
      { method: "POST", encType: "application/json" }
    )
    setEditing(false)
  }

  const hasAddress = profile.street || profile.city

  if (!editing) {
    return (
      <section className="rounded-xl border bg-card p-6 space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-base font-semibold">{t("account.sections.address.title")}</h2>
          <Button variant="outline" size="sm" onClick={handleEdit}>
            {t("account.edit")}
          </Button>
        </div>
        {hasAddress && (
          <dl className="space-y-3 text-sm">
            <div>
              <dt className="text-muted-foreground">{t("account.sections.address.street")}</dt>
              <dd className="font-medium mt-0.5">{profile.street || "—"}</dd>
            </div>
            {profile.complement && (
              <div>
                <dt className="text-muted-foreground">{t("account.sections.address.complement")}</dt>
                <dd className="font-medium mt-0.5">{profile.complement}</dd>
              </div>
            )}
            <div>
              <dt className="text-muted-foreground">{t("account.sections.address.neighborhood")}</dt>
              <dd className="font-medium mt-0.5">{profile.neighborhood || "—"}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">{t("account.sections.address.city")}</dt>
              <dd className="font-medium mt-0.5">{[profile.city, profile.state].filter(Boolean).join(", ") || "—"}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">{t("account.sections.address.zip")}</dt>
              <dd className="font-medium mt-0.5">{profile.zip || "—"}</dd>
            </div>
          </dl>
        )}
        {!hasAddress && (
          <p className="text-sm text-muted-foreground">—</p>
        )}
        {fetcher.data && "error" in fetcher.data && (
          <p className="text-xs text-destructive">{t("account.saveError")}</p>
        )}
        {fetcher.data && "ok" in fetcher.data && fetcher.data.intent === "save-address" && (
          <p className="text-xs text-green-600">{t("account.saveSuccess")}</p>
        )}
      </section>
    )
  }

  return (
    <section className="rounded-xl border bg-card p-6 space-y-4">
      <h2 className="text-base font-semibold">{t("account.sections.address.title")}</h2>

      {hasActiveSubscription && (
        <div className="flex gap-2 rounded-lg border border-yellow-300 bg-yellow-50 px-3 py-2 text-sm text-yellow-800">
          <AlertTriangle className="mt-0.5 shrink-0 size-4" />
          <span>{t("account.sections.address.activeSubscriptionWarning")}</span>
        </div>
      )}

      <div className="space-y-3">
        <div className="space-y-1">
          <label className="text-sm text-muted-foreground">{t("account.sections.address.street")}</label>
          <Input value={fields.street} onChange={set("street")} />
          {errors.street && <p className="text-xs text-destructive">{errors.street}</p>}
        </div>
        <div className="space-y-1">
          <label className="text-sm text-muted-foreground">{t("account.sections.address.complement")}</label>
          <Input value={fields.complement} onChange={set("complement")} />
        </div>
        <div className="grid grid-cols-2 gap-3">
          <div className="space-y-1">
            <label className="text-sm text-muted-foreground">{t("account.sections.address.neighborhood")}</label>
            <Input value={fields.neighborhood} onChange={set("neighborhood")} />
            {errors.neighborhood && <p className="text-xs text-destructive">{errors.neighborhood}</p>}
          </div>
          <div className="space-y-1">
            <label className="text-sm text-muted-foreground">{t("account.sections.address.zip")}</label>
            <Input value={fields.zip} onChange={set("zip")} placeholder={ZIP_PLACEHOLDER} />
            {errors.zip && <p className="text-xs text-destructive">{errors.zip}</p>}
          </div>
        </div>
        <div className="grid grid-cols-2 gap-3">
          <div className="space-y-1">
            <label className="text-sm text-muted-foreground">{t("account.sections.address.city")}</label>
            <Input value={fields.city} onChange={set("city")} />
            {errors.city && <p className="text-xs text-destructive">{errors.city}</p>}
          </div>
          <div className="space-y-1">
            <label className="text-sm text-muted-foreground">{t("account.sections.address.state")}</label>
            <Input value={fields.state} onChange={set("state")} maxLength={2} className="uppercase" />
            {errors.state && <p className="text-xs text-destructive">{errors.state}</p>}
          </div>
        </div>
      </div>
      <div className="flex gap-2 justify-end">
        <Button variant="ghost" size="sm" onClick={handleCancel} disabled={isSaving}>
          {t("account.cancel")}
        </Button>
        <Button size="sm" onClick={handleSave} disabled={isSaving}>
          {isSaving ? t("account.saving") : t("account.save")}
        </Button>
      </div>
    </section>
  )
}
