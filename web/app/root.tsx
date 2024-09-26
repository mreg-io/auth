import {
  isRouteErrorResponse,
  Links,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
  useRouteError,
} from "@remix-run/react";
import { LinksFunction } from "@remix-run/node";
import { ReactNode } from "react";
import { Internal, NotFound } from "~/components/errors";
import styles from "./tailwind.css?url";

export const links: LinksFunction = () => [{ rel: "stylesheet", href: styles }];

type LayoutProps = { children: ReactNode; title?: string };

function Layout({ children, title }: LayoutProps) {
  return (
    <html lang="en">
      <head>
        {title ? <title>{title}</title> : undefined}
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <Meta />
        <Links />
      </head>
      <body>
        {children}
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
}

export default function App() {
  return (
    <Layout>
      <Outlet />
    </Layout>
  );
}

export function ErrorBoundary() {
  const error = useRouteError();
  if (isRouteErrorResponse(error)) {
    switch (error.status) {
      case 404:
        return (
          <Layout title="Not Found | My Registry">
            <NotFound />
          </Layout>
        );
      default:
        return (
          <Layout title="Internal Server Error | My Registry">
            <Internal status={error.status} />
          </Layout>
        );
    }
  } else {
    return <h1 className="text-lg">Internal Server Error</h1>;
  }
}
