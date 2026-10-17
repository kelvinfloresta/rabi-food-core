import { redirect } from "react-router"
import type { Route } from "./+types/sign-out"

export async function action(_: Route.ActionArgs) {
  return redirect("/", {
    headers: {
      "Set-Cookie": "auth_token=; Path=/; HttpOnly; SameSite=Lax; Max-Age=0",
    },
  })
}
