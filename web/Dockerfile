FROM --platform=$BUILDPLATFORM node:22 AS base

FROM --platform=$BUILDPLATFORM base AS deps

WORKDIR /app

COPY package.json yarn.lock* package-lock.json* pnpm-lock.yaml* ./

RUN --mount=type=secret,id=npmrc,target=/root/.npmrc \
  if [ -f yarn.lock ]; then yarn --frozen-lockfile; \
  elif [ -f package-lock.json ]; then npm ci; \
  elif [ -f pnpm-lock.yaml ]; then corepack enable pnpm && pnpm i --frozen-lockfile; \
  else echo "Lockfile not found." && exit 1; \
  fi

FROM --platform=$BUILDPLATFORM base AS prod-deps

WORKDIR /app

COPY package.json yarn.lock* package-lock.json* pnpm-lock.yaml* ./

RUN --mount=type=secret,id=npmrc,target=/root/.npmrc \
  if [ -f yarn.lock ]; then yarn --frozen-lockfile; \
  elif [ -f package-lock.json ]; then npm ci; \
  elif [ -f pnpm-lock.yaml ]; then corepack enable pnpm && pnpm i -P --frozen-lockfile; \
  else echo "Lockfile not found." && exit 1; \
  fi

FROM --platform=$BUILDPLATFORM base AS builder

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules

COPY . .

RUN \
  if [ -f yarn.lock ]; then yarn run build; \
  elif [ -f package-lock.json ]; then npm run build; \
  elif [ -f pnpm-lock.yaml ]; then corepack enable pnpm && pnpm run build; \
  else echo "Lockfile not found." && exit 1; \
  fi

FROM gcr.io/distroless/nodejs22-debian12:nonroot

WORKDIR /app

ENV NODE_ENV=production

COPY --from=prod-deps --chown=nonroot:nonroot /app/package.json /app/pnpm-lock.yaml ./

COPY --from=prod-deps --chown=nonroot:nonroot /app/node_modules ./node_modules

COPY --from=builder --chown=nonroot:nonroot /app/build/ ./build

ENV PORT=3000

ENV HOSTNAME="0.0.0.0"

EXPOSE 3000

CMD ["./node_modules/@remix-run/serve/dist/cli.js", "./build/server/index.js"]
