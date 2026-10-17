import { resolve } from "node:path"
import { createInstance, type i18n } from "i18next"
import { createCookie, type RouterContextProvider } from "react-router"
import { createI18nextMiddleware } from "remix-i18next"
import FsBackend from "i18next-fs-backend"
import { initReactI18next } from "react-i18next"
import i18nConfig from "~/i18n"

export const localeCookie = createCookie("i18n", {
  path: "/",
  sameSite: "lax",
  maxAge: 60 * 60 * 24 * 365,
})

const serverI18nConfig = {
  ...i18nConfig,
  backend: {
    loadPath: resolve("./public/locales/{{lng}}/{{ns}}.json"),
  },
}

export const [i18nextMiddleware, getLocale, getInstance] = createI18nextMiddleware({
  detection: {
    supportedLanguages: [...i18nConfig.supportedLngs],
    fallbackLanguage: i18nConfig.fallbackLng,
    order: ["cookie", "header"],
    cookie: localeCookie,
  },
  i18next: serverI18nConfig,
  plugins: [FsBackend, initReactI18next],
})

async function createFallbackInstance(): Promise<i18n> {
  const instance = createInstance(serverI18nConfig)
  instance.use(FsBackend).use(initReactI18next)
  await instance.init({ lng: i18nConfig.fallbackLng })
  return instance
}

// Unmatched URLs (404) skip route middleware, so there is no request-scoped instance to render with
export async function getRenderInstance(context: Readonly<RouterContextProvider>): Promise<i18n> {
  try {
    return getInstance(context)
  } catch {
    return createFallbackInstance()
  }
}
