import { Form, Link, useNavigation } from "@remix-run/react";
import { cn } from "~/lib/utils";
import { Button, buttonVariants } from "~/components/ui/button";
import { Label } from "~/components/ui/label";
import { Input } from "~/components/ui/input";
import { Icons } from "~/components/icons";
import {
  ActionFunctionArgs,
  LoaderFunctionArgs,
  MetaFunction,
} from "@remix-run/node";
import { registrationService } from "~/lib/connect.server";
import { useEffect, useState } from "react";
import { parseCookie } from "~/lib/cookie";
import { CreateRegistrationFlowResponse } from "@buf/mreg_protobuf.bufbuild_es/mreg/auth/v1alpha1/registration_service_pb";
import { protobuf, useLoaderProtobuf } from "~/lib/protobuf";

export const meta: MetaFunction = () => [
  { title: "Create an Account | My Registry" },
];

export async function loader({ request }: LoaderFunctionArgs) {
  const { response, headers } =
    await registrationService.createRegistrationFlow(
      {},
      {
        headers: {
          "X-Forwarded-For": "0.0.0.0",
          "User-Agent": request.headers.get("User-Agent")!,
        },
      }
    );

  return protobuf(response, { headers });
}

export async function action({ request }: ActionFunctionArgs) {
  const data = Object.fromEntries(await request.formData());
  console.log(data);
  return null;
}

export default function Registration() {
  const flow = useLoaderProtobuf(CreateRegistrationFlowResponse);
  const { formAction } = useNavigation();
  const isSubmitting = formAction === "/registration";

  const [csrfToken, setCsrfToken] = useState<string>();
  useEffect(() => {
    const cookies = parseCookie();
    setCsrfToken(cookies.get("csrf_token"));
  }, []);

  return (
    <>
      <div className="container relative h-full flex-col items-center justify-center grid md:max-w-none lg:grid-cols-2 lg:px-0">
        <Link
          to="/login"
          className={cn(
            buttonVariants({ variant: "ghost" }),
            "absolute right-4 top-4 md:right-8 md:top-8"
          )}
        >
          Login
        </Link>
        <div className="relative hidden h-full flex-col bg-muted p-10 text-white dark:border-r lg:flex">
          <div className="absolute inset-0 bg-zinc-900" />
          <div className="relative z-20 flex items-center text-lg font-medium">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
              className="mr-2 h-6 w-6"
            >
              <path d="M15 6v12a3 3 0 1 0 3-3H6a3 3 0 1 0 3 3V6a3 3 0 1 0-3 3h12a3 3 0 1 0-3-3" />
            </svg>
            My Registry
          </div>
          <div className="relative z-20 mt-auto">
            <blockquote className="space-y-2">
              <p className="text-lg">
                &ldquo;This library has saved me countless hours of work and
                helped me deliver stunning designs to my clients faster than
                ever before.&rdquo;
              </p>
              <footer className="text-sm">Sofia Davis</footer>
            </blockquote>
          </div>
        </div>
        <div className="p-8">
          <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
            <div className="flex flex-col space-y-2 text-center">
              <h1 className="text-2xl font-semibold tracking-tight">
                Create an account
              </h1>
              <p className="text-sm text-muted-foreground">
                Enter your email below to create your account
              </p>
            </div>
            <div className={cn("grid gap-6")}>
              <Form method="POST" action="/registration">
                <fieldset className="grid gap-4" disabled={isSubmitting}>
                  <div className="grid gap-2">
                    <Label htmlFor="email">Email</Label>
                    <Input
                      id="email"
                      name="email"
                      placeholder="name@example.com"
                      type="email"
                      required
                      autoCapitalize="none"
                      autoComplete="email"
                      autoCorrect="off"
                    />
                  </div>
                  <div className="grid gap-2">
                    <Label htmlFor="password">Password</Label>
                    <Input
                      id="password"
                      name="password"
                      type="password"
                      required
                      minLength={8}
                      maxLength={256}
                      autoCapitalize="none"
                      autoComplete="new-password"
                      autoCorrect="off"
                    />
                  </div>
                  <input
                    name="flow-name"
                    className="hidden"
                    defaultValue={flow.registrationFlow?.name}
                  />
                  <input
                    name="flow-etag"
                    className="hidden"
                    defaultValue={flow.registrationFlow?.etag}
                  />
                  <input
                    name="csrf-token"
                    className="hidden"
                    defaultValue={csrfToken}
                  />
                  <Button type="submit">
                    {isSubmitting && (
                      <Icons.spinner className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    Sign Up
                  </Button>
                </fieldset>
              </Form>
            </div>
            <p className="px-8 text-center text-sm text-muted-foreground">
              By clicking continue, you agree to our{" "}
              <Link
                to="/terms"
                className="underline underline-offset-4 hover:text-primary"
              >
                Terms of Service
              </Link>{" "}
              and{" "}
              <Link
                to="/privacy"
                className="underline underline-offset-4 hover:text-primary"
              >
                Privacy Policy
              </Link>
              .
            </p>
          </div>
        </div>
      </div>
    </>
  );
}
