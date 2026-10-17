import { redirect } from "react-router"
import type { Route } from "./+types/sign-out"
import { apiClient } from "~/lib/api.server"
import { destroySessionHeaders, getAccessToken } from "~/lib/session.server"

export async function action({ context }: Route.ActionArgs) {
  const token = getAccessToken(context)

  // Revoke the refresh token on the auth server; cookies are cleared even if this fails
  if (token) await apiClient().post("/auth/signout", null, token)

  return redirect("/", { headers: await destroySessionHeaders() })
}
