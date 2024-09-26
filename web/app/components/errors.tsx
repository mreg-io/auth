import { Link } from "@remix-run/react";
import { Button } from "~/components/shadcn/button";

export function NotFound() {
  return (
    <main className="grid min-h-full place-items-center bg-white px-6 py-24 sm:py-32 lg:px-8">
      <div className="text-center">
        <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
          404
        </h1>
        <p className="mt-6 text-base text-muted-foreground">
          The page doesn&apos;t exist.
        </p>
        <Button variant="default" size="lg" className="mt-8" asChild>
          <Link to="/login">Log in to My Registry</Link>
        </Button>
      </div>
    </main>
  );
}

export type InternalProps = { status?: number };

export function Internal({ status }: InternalProps) {
  return (
    <main className="grid min-h-full place-items-center bg-white px-6 py-24 sm:py-32 lg:px-8">
      <div className="text-center">
        <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl">
          {status ?? 500}
        </h1>
        <p className="mt-6 text-base text-muted-foreground">
          Oops, something went wrong.
        </p>
        <Button variant="default" size="lg" className="mt-8" asChild>
          <Link to="/login">Log in to My Registry</Link>
        </Button>
      </div>
    </main>
  );
}
