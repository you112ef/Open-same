# Cloudflare Pages Chat (OpenAI)

This is a minimal Cloudflare Pages project with:

- Static UI in `public/`
- Pages Function at `functions/api/chat.js` that calls OpenAI Chat Completions

## Local development

1. Install deps:

```bash
cd cf-pages-app
npm install
```

2. Export your OpenAI key for local dev (or set in `wrangler.toml` [vars] temporarily):

```bash
export OPENAI_API_KEY=sk-... # do NOT commit
```

3. Start dev server:

```bash
npm run dev
```

Open http://localhost:8788 and chat.

## Deploy to Cloudflare Pages

1. Create a Pages project in the Cloudflare dashboard with this repo/folder.
2. Set a Pages secret `OPENAI_API_KEY` in Project Settings → Environment Variables → Production/Preview.
3. Deploy via CLI:

```bash
npm run deploy
```

Wrangler will upload `public/` and include functions from `functions/`.
