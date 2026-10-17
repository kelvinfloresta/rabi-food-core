import {
  createContext,
  createCookie,
  type Cookie,
  type MiddlewareFunction,
  type RouterContextProvider,
} from "react-router"
import { apiClient } from "~/lib/api.server"

const ACCESS_TOKEN_MAX_AGE_SECONDS = 60 * 60
const REFRESH_TOKEN_MAX_AGE_SECONDS = 60 * 60 * 24 * 30
// Renew slightly before expiry so the token does not expire mid-request on its way to the API
const EXPIRY_LEEWAY_SECONDS = 30
const MILLISECONDS_PER_SECOND = 1000
const EXPIRED_COOKIE_MAX_AGE = 0

export type TokenPair = Readonly<{
  accessToken: string
  refreshToken: string
}>

type JwtPayload = Readonly<{ exp?: number }>

const baseCookieOptions = {
  path: "/",
  httpOnly: true,
  sameSite: "lax",
  secure: process.env.NODE_ENV === "production",
} as const

const accessTokenCookie = createCookie("auth_token", { ...baseCookieOptions, maxAge: ACCESS_TOKEN_MAX_AGE_SECONDS })
const refreshTokenCookie = createCookie("refresh_token", { ...baseCookieOptions, maxAge: REFRESH_TOKEN_MAX_AGE_SECONDS })

const accessTokenContext = createContext<string | null>(null)

function isTokenExpired(token: string): boolean {
  try {
    const encodedPayload = token.split(".")[1] ?? ""
    const payload = JSON.parse(Buffer.from(encodedPayload, "base64url").toString()) as JwtPayload
    if (typeof payload.exp !== "number") return true
    return Date.now() >= (payload.exp - EXPIRY_LEEWAY_SECONDS) * MILLISECONDS_PER_SECOND
  } catch {
    return true
  }
}

async function parseToken(cookie: Cookie, cookieHeader: string | null): Promise<string | null> {
  const value: unknown = await cookie.parse(cookieHeader)
  if (typeof value !== "string" || !value) return null
  return value
}

function getCookieName(setCookie: string): string {
  return setCookie.split("=")[0] ?? ""
}

// A cookie set by the loader/action itself (sign-in, sign-out) wins over the one the middleware would add
function appendMissingCookies(response: Response, headers: Headers): Response {
  const alreadySet = new Set(response.headers.getSetCookie().map(getCookieName))
  for (const setCookie of headers.getSetCookie()) {
    if (alreadySet.has(getCookieName(setCookie))) continue
    response.headers.append("Set-Cookie", setCookie)
  }
  return response
}

async function refreshSession(refreshToken: string): Promise<TokenPair | null> {
  const result = await apiClient().post<TokenPair>("/auth/refresh", { refreshToken })
  if (!result.ok) return null
  return result.data
}

export async function createSessionHeaders(tokens: TokenPair): Promise<Headers> {
  const headers = new Headers()
  headers.append("Set-Cookie", await accessTokenCookie.serialize(tokens.accessToken))
  headers.append("Set-Cookie", await refreshTokenCookie.serialize(tokens.refreshToken))
  return headers
}

export async function destroySessionHeaders(): Promise<Headers> {
  const headers = new Headers()
  headers.append("Set-Cookie", await accessTokenCookie.serialize("", { maxAge: EXPIRED_COOKIE_MAX_AGE }))
  headers.append("Set-Cookie", await refreshTokenCookie.serialize("", { maxAge: EXPIRED_COOKIE_MAX_AGE }))
  return headers
}

export function getAccessToken(context: Readonly<RouterContextProvider>): string | null {
  return context.get(accessTokenContext)
}

// Auth tokens never reach the browser: it only holds HttpOnly cookies, and an expired
// access token is renewed here before any loader/action runs
export const authMiddleware: MiddlewareFunction<Response> = async ({ request, context }, next) => {
  const cookieHeader = request.headers.get("Cookie")

  const accessToken = await parseToken(accessTokenCookie, cookieHeader)
  if (accessToken && !isTokenExpired(accessToken)) {
    context.set(accessTokenContext, accessToken)
    return next()
  }

  const refreshToken = await parseToken(refreshTokenCookie, cookieHeader)
  if (!refreshToken) return next()

  const tokens = await refreshSession(refreshToken)
  if (!tokens) return appendMissingCookies(await next(), await destroySessionHeaders())

  context.set(accessTokenContext, tokens.accessToken)
  return appendMissingCookies(await next(), await createSessionHeaders(tokens))
}
