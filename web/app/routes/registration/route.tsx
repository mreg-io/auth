import { useNavigate, useNavigation } from "@remix-run/react";
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
import RegistrationForm from "~/routes/registration/registration-form";
import { useToast } from "~/hooks/use-toast";

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
    <RegistrationForm
      loading={isSubmitting}
      disabled={isSubmitting}
      csrfToken={csrfToken}
    >
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
