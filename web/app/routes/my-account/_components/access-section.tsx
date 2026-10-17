import { useTranslation } from "react-i18next"
import { Button } from "~/components/ui/button"
import { Lock } from "lucide-react"

type AccessSectionProps = Readonly<{
  email: string | undefined
}>

export function AccessSection({ email }: AccessSectionProps) {
  const { t } = useTranslation()

  return (
    <section className="rounded-xl border bg-card p-6 space-y-4">
      <h2 className="text-base font-semibold">{t("account.sections.access.title")}</h2>
      <dl className="space-y-3 text-sm">
        <div>
          <dt className="text-muted-foreground">{t("account.sections.access.email")}</dt>
          <dd className="font-medium mt-0.5">{email ?? "—"}</dd>
          <p className="text-xs text-muted-foreground mt-1">{t("account.sections.access.emailReadOnlyHint")}</p>
        </div>
      </dl>
      <div className="pt-1">
        <Button variant="outline" size="sm" disabled>
          <Lock />
          {t("account.sections.access.changePassword")}
        </Button>
      </div>
    </section>
  )
}
