import { useState } from "react"
import { useFetcher } from "react-router"
import { useTranslation } from "react-i18next"
import PhoneInput, { isValidPhoneNumber } from "react-phone-number-input"
import "react-phone-number-input/style.css"
import { Button } from "~/components/ui/button"
import { Input } from "~/components/ui/input"
import type { AccountProfile } from "~/components/subscription/types"
import type { AccountActionResult } from "../index"

type PersonalInfoSectionProps = Readonly<{
  profile: AccountProfile
}>

type PersonalInfoErrors = Readonly<{
  name?: string
  phone?: string
}>

type MutablePersonalInfoErrors = {
  name?: string
  phone?: string
}

const PHONE_PLACEHOLDER = "+55 (11) 99999-9999"

export function PersonalInfoSection({ profile }: PersonalInfoSectionProps) {
  const { t } = useTranslation()
  const fetcher = useFetcher<AccountActionResult>()
  const [editing, setEditing] = useState(false)
  const [nameValue, setNameValue] = useState(profile.name ?? "")
  const [phoneValue, setPhoneValue] = useState(profile.phone ?? "")
  const [errors, setErrors] = useState<PersonalInfoErrors>({})

  const isSaving = fetcher.state !== "idle"

  function handleEdit() {
    setNameValue(profile.name ?? "")
    setPhoneValue(profile.phone ?? "")
    setErrors({})
    setEditing(true)
  }

  function handleCancel() {
    setEditing(false)
    setErrors({})
  }

  function validate(): boolean {
    const next: MutablePersonalInfoErrors = {}
    if (!nameValue.trim()) next.name = t("account.sections.personalInfo.nameRequired")
    if (!phoneValue || !isValidPhoneNumber(phoneValue)) next.phone = t("account.sections.personalInfo.phoneInvalid")
    setErrors(next)
    return Object.keys(next).length === 0
  }

  function handleSave() {
    if (!validate()) return
    fetcher.submit(
      { intent: "save-personal-info", name: nameValue, phone: phoneValue },
      { method: "POST", encType: "application/json" }
    )
    setEditing(false)
  }

  if (!editing) {
    return (
      <section className="rounded-xl border bg-card p-6 space-y-4">
        <div className="flex items-center justify-between">
          <h2 className="text-base font-semibold">{t("account.sections.personalInfo.title")}</h2>
          <Button variant="outline" size="sm" onClick={handleEdit}>
            {t("account.edit")}
          </Button>
        </div>
        <dl className="space-y-3 text-sm">
          <div>
            <dt className="text-muted-foreground">{t("account.sections.personalInfo.name")}</dt>
            <dd className="font-medium mt-0.5">{profile.name ?? "—"}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">{t("account.sections.personalInfo.taxId")}</dt>
            <dd className="font-medium mt-0.5">{profile.taxId ?? "—"}</dd>
          </div>
          <div>
            <dt className="text-muted-foreground">{t("account.sections.personalInfo.phone")}</dt>
            <dd className="font-medium mt-0.5">{profile.phone || "—"}</dd>
          </div>
        </dl>
        {fetcher.data && "error" in fetcher.data && (
          <p className="text-xs text-destructive">{t("account.saveError")}</p>
        )}
        {fetcher.data && "ok" in fetcher.data && fetcher.data.intent === "save-personal-info" && (
          <p className="text-xs text-green-600">{t("account.saveSuccess")}</p>
        )}
      </section>
    )
  }

  return (
    <section className="rounded-xl border bg-card p-6 space-y-4">
      <h2 className="text-base font-semibold">{t("account.sections.personalInfo.title")}</h2>
      <div className="space-y-3">
        <div className="space-y-1">
          <label className="text-sm text-muted-foreground">{t("account.sections.personalInfo.name")}</label>
          <Input value={nameValue} onChange={(e) => setNameValue(e.target.value)} />
          {errors.name && <p className="text-xs text-destructive">{errors.name}</p>}
        </div>
        <div className="space-y-1">
          <label className="text-sm text-muted-foreground">{t("account.sections.personalInfo.phone")}</label>
          <PhoneInput
            international
            defaultCountry="BR"
            value={phoneValue}
            onChange={(v) => setPhoneValue(v ?? "")}
            placeholder={PHONE_PLACEHOLDER}
            className="flex h-9 w-full rounded-lg border border-input bg-transparent px-3 py-1 text-sm shadow-xs transition-colors placeholder:text-muted-foreground focus-within:ring-2 focus-within:ring-ring/50 focus-within:border-ring"
          />
          {errors.phone && <p className="text-xs text-destructive">{errors.phone}</p>}
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
