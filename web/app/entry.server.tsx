import { createReadableStreamFromReadable } from "@react-router/node"
import { PassThrough } from "node:stream"
import { renderToPipeableStream } from "react-dom/server"
import { I18nextProvider } from "react-i18next"
import type { EntryContext, RouterContextProvider } from "react-router"
import { ServerRouter } from "react-router"
import { getRenderInstance } from "~/i18n.server"

const ABORT_DELAY = 5000

export default async function handleRequest(
  request: Request,
  responseStatusCode: number,
  responseHeaders: Headers,
  routerContext: EntryContext,
  loadContext: RouterContextProvider
) {
  const instance = await getRenderInstance(loadContext)

  return new Promise<Response>((resolve, reject) => {
    let shellRendered = false

    const { pipe, abort } = renderToPipeableStream(
      <I18nextProvider i18n={instance}>
        <ServerRouter context={routerContext} url={request.url} />
      </I18nextProvider>,
      {
        onShellReady() {
          shellRendered = true
          const body = new PassThrough()
          const stream = createReadableStreamFromReadable(body)
          responseHeaders.set("Content-Type", "text/html")
          resolve(
            new Response(stream, {
              headers: responseHeaders,
              status: responseStatusCode,
            })
          )
          pipe(body)
        },
        onShellError(error: unknown) {
          reject(error)
        },
        onError(error: unknown) {
          responseStatusCode = 500
          if (shellRendered) console.error(error)
        },
      }
    )

    setTimeout(abort, ABORT_DELAY)
  })
}
