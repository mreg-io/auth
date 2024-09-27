declare global {
  namespace NodeJS {
    interface ProcessEnv {
      AUTH_API_URL: string;
      CSRF_TOKEN_SECRET: string;
    }
  }
}

export {};
