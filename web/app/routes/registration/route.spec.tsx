import { describe, expect, it } from "vitest";
import { createRemixStub } from "@remix-run/testing";
import Registration from "~/routes/registration/route";
import { render, screen } from "@testing-library/react";
import { unstable_data as data } from "@remix-run/react";

describe("Registration", () => {
  it("should render registration form", async () => {
    const RemixStub = createRemixStub([
      {
        path: "/registration",
        Component: Registration,
        loader: () => {
          return data({
            response: {
              registrationFlow: {
                name: "registrationFlows/01923634-d98c-9563-8c9e-3a676d49ac00",
                flowId: "01923634-d98c-9563-8c9e-3a676d49ac00",
                etag: "UWu5u//dU1PuukmRaEmO1RUNSN5NkUgXV/3gpUhMHow=.MDE5MjM2MzQtZDk4My05MjM1LTU2OWMtN2E1ZTc2NjcxOWUzIWJmMDFiODA3LWIwZTQtNDkxOS1hZGFiLTM5NTQzZDU5OTgwYg==",
              },
            },
            csrfToken:
              "+oQ30hU26lls4/8uXsQCzyCLaI6WJ/JzzljgkADdhuc=.NzQ1OGFmNzItZWE2My00YmNkLTlhNzAtOTA5Mjc3M2NkZWY1ITljMDQ0MGFjLTZkZTgtNDdmNS05MDJiLTRjZjI5YTQ1ZjJmNg==",
          });
        },
      },
    ]);
    render(<RemixStub initialEntries={["/registration"]} />);

    const heading = await screen.findByRole("heading");
    expect(heading).toHaveTextContent("Create an account");

    const subtitle = await screen.findByText(
      "Enter your email below to create your account",
    );
    expect(subtitle).toBeVisible();

    const email = await screen.findByLabelText("Email");
    expect(email).toBeVisible();
    expect(email).toHaveAttribute("type", "email");
    expect(email).toHaveAttribute("required");
    expect(email).toHaveAttribute("placeholder", "name@example.com");
    expect(email).toHaveAttribute("autocomplete", "email");

    const password = await screen.findByLabelText("Password");
    expect(password).toBeVisible();
    expect(password).toHaveAttribute("type", "password");
    expect(password).toHaveAttribute("required");
    expect(password).toHaveAttribute("minLength", "8");
    expect(password).toHaveAttribute("maxLength", "256");
    expect(password).toHaveAttribute("autocomplete", "new-password");

    const signUp = await screen.findByRole("button");
    expect(signUp).toBeVisible();
    expect(signUp).toHaveTextContent("Sign Up");
    expect(signUp).toHaveAttribute("type", "submit");

    const login = await screen.findByRole("link", { name: "Login" });
    expect(login).toBeVisible();
    expect(login).toHaveAttribute("href", "/login");

    const terms = await screen.findByRole("link", { name: "Terms of Service" });
    expect(terms).toBeVisible();
    expect(terms).toHaveAttribute("href", "/terms");

    const privacy = await screen.findByRole("link", { name: "Privacy Policy" });
    expect(privacy).toBeVisible();
    expect(privacy).toHaveAttribute("href", "/privacy");
  });
});
