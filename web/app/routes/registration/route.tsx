import { useLoaderData, useNavigate, useNavigation } from "@remix-run/react";
import {
  unstable_data as data,
  ActionFunctionArgs,
  LoaderFunctionArgs,
  MetaFunction,
} from "@remix-run/node";
import { registrationService } from "~/lib/connect.server";
import { useEffect } from "react";
import RegistrationForm from "~/routes/registration/registration-form";
import { useToast } from "~/hooks/use-toast";
import { generateCSRFTokenFromHeaders } from "~/lib/csrf.server";

export const meta: MetaFunction = () => [
  { title: "Create an Account | My Registry" },
];

export const loader = async ({ request }: LoaderFunctionArgs) => {
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
  const csrfToken = generateCSRFTokenFromHeaders(headers);

  return data(
    {
      response,
      csrfToken,
    },
    { headers }
  );
};

export async function action({ request }: ActionFunctionArgs) {
  const data = Object.fromEntries(await request.formData());
  console.log(data);
  return null;
}

export default function Registration() {
  const { response, csrfToken } = useLoaderData<typeof loader>();
  const { formAction } = useNavigation();
  const isSubmitting = formAction === "/registration";

  return (
    <RegistrationForm loading={isSubmitting} disabled={isSubmitting}>
      <input
        name="flow-name"
        className="hidden"
        defaultValue={response.registrationFlow?.name}
      />
      <input
        name="flow-etag"
        className="hidden"
        defaultValue={response.registrationFlow?.etag}
      />
      <input name="csrf-token" className="hidden" defaultValue={csrfToken} />
    </RegistrationForm>
  );
}

export function ErrorBoundary() {
  const navigate = useNavigate();
  const { toast } = useToast();

  useEffect(() => {
    toast({
      variant: "destructive",
      title: "Uh oh! Something went wrong.",
      description: "There was a problem with your request.",
    });
  }, [toast, navigate]);

  return <RegistrationForm disabled />;
}
