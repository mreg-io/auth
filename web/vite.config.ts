import { vitePlugin as remix } from "@remix-run/dev";
import { defineConfig, loadEnv } from "vite";
import { coverageConfigDefaults } from "vitest/config";
import tsconfigPaths from "vite-tsconfig-paths";
import * as process from "node:process";

export default defineConfig(({ mode }) => ({
  plugins: [
    !process.env.VITEST &&
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
    environmentMatchGlobs: [
      ["**/*.[jt]sx", "jsdom"],
      ["**/*", "node"],
    ],
    env: loadEnv(mode, process.cwd(), "") as Partial<NodeJS.ProcessEnv>,
    setupFiles: "vitest-setup.ts",
    restoreMocks: true,
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
}));

declare module "@remix-run/node" {
  // or cloudflare, deno, etc.
  interface Future {
    unstable_singleFetch: true;
  }
}
