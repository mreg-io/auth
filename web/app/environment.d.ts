declare global {
  namespace NodeJS {
    interface ProcessEnv {
      AUTH_API_URL: string;
    }
  }
}

export {};
