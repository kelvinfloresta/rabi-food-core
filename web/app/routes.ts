import { type RouteConfig, index, layout, route } from "@react-router/dev/routes";


export default [
  route(".well-known/*", "routes/well-known.tsx"),
  route("sign-out", "routes/sign-out.tsx"),
  route("my-account", "routes/my-account/index.tsx"),
  route("my-subscriptions", "routes/my-subscriptions/index.tsx"),
  layout("routes/storefront/layout.tsx", [
    index("routes/storefront/storefront-page.tsx"),
    route("sobre-assinatura", "routes/storefront/apresentacao.tsx"),
    route("product/:id", "routes/storefront/product.tsx"),
    route("subscribe", "routes/storefront/subscribe.tsx"),
  ]),
] satisfies RouteConfig;
