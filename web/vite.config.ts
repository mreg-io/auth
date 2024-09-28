import { vitePlugin as remix } from "@remix-run/dev";
import { defineConfig } from "vite";
import { coverageConfigDefaults } from "vitest/config";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig({
  plugins: [
    remix({
      future: {
        v3_fetcherPersist: true,
        v3_relativeSplatPath: true,
        v3_throwAbortReason: true,
        unstable_singleFetch: true,
      },
    }),
    tsconfigPaths(),
  ],
  test: {
    typecheck: {
      enabled: true,
    },
    coverage: {
      provider: "istanbul",
      exclude: [
        ...coverageConfigDefaults.exclude,
        "build/**",
        // Shadcn generated code
        "app/components/shadcn/**",
        "app/hooks/use-toast.ts",
      ],
      reporter: ["text", "html", "cobertura"],
    },
    reporters: ["default", ["junit", { suiteName: "web unit test" }]],
    outputFile: "junit.xml",
  },
});

declare module "@remix-run/node" {
  // or cloudflare, deno, etc.
  interface Future {
    unstable_singleFetch: true;
  }
}
