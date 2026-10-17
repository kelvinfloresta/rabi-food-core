import { NavLink } from "react-router"
import { useTranslation } from "react-i18next"
import { Home, Package, User } from "lucide-react"
import { cn } from "~/lib/utils"

type NavItem = Readonly<{
  to: string
  icon: React.ComponentType<{ className?: string }>
  labelKey: string
  end: boolean
}>

const NAV_ITEMS: readonly NavItem[] = [
  { to: "/", icon: Home, labelKey: "nav.storefront", end: true },
  { to: "/my-subscriptions", icon: Package, labelKey: "nav.subscriptions", end: false },
  { to: "/my-account", icon: User, labelKey: "nav.account", end: false },
] as const

export const BOTTOM_NAV_HEIGHT = 64

export function BottomNav() {
  const { t } = useTranslation()

  return (
    <nav
      className="fixed bottom-0 left-0 right-0 z-50 border-t bg-background"
      style={{ paddingBottom: "env(safe-area-inset-bottom)" }}
    >
      <div className="flex h-16">
        {NAV_ITEMS.map(({ to, icon: Icon, labelKey, end }) => (
          <NavLink
            key={to}
            to={to}
            end={end}
            className={({ isActive }) =>
              cn(
                "flex flex-1 flex-col items-center justify-center gap-1 text-xs transition-colors",
                isActive ? "text-primary" : "text-muted-foreground hover:text-foreground"
              )
            }
          >
            <Icon className="size-5" />
            <span>{t(labelKey)}</span>
          </NavLink>
        ))}
      </div>
    </nav>
  )
}
