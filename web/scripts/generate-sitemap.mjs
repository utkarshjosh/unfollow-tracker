import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const webRoot = path.resolve(__dirname, '..');
const appFile = path.join(webRoot, 'src', 'App.tsx');
const outputFile = path.join(webRoot, 'public', 'sitemap.xml');

const DEFAULT_SITE_URL = 'https://unfollowtracker.com';
const EXCLUDED_ROUTES = new Set([
  '/dashboard',
]);

function loadEnvFile(filePath) {
  if (!fs.existsSync(filePath)) return;
  const content = fs.readFileSync(filePath, 'utf-8');
  for (const rawLine of content.split('\n')) {
    const line = rawLine.trim();
    if (!line || line.startsWith('#')) continue;
    const equalIndex = line.indexOf('=');
    if (equalIndex <= 0) continue;
    const key = line.slice(0, equalIndex).trim();
    const value = line.slice(equalIndex + 1).trim();
    if (!(key in process.env)) {
      process.env[key] = value.replace(/^['"]|['"]$/g, '');
    }
  }
}

function loadBuildEnv() {
  const mode = process.env.NODE_ENV === 'development' ? 'development' : 'production';
  loadEnvFile(path.join(webRoot, '.env'));
  loadEnvFile(path.join(webRoot, '.env.local'));
  loadEnvFile(path.join(webRoot, `.env.${mode}`));
  loadEnvFile(path.join(webRoot, `.env.${mode}.local`));
}

function normalizeSiteUrl(value) {
  const trimmed = (value || DEFAULT_SITE_URL).trim();
  const withoutTrailingSlash = trimmed.replace(/\/+$/, '');
  return withoutTrailingSlash || DEFAULT_SITE_URL;
}

function routePriority(route) {
  if (route === '/') return '1.0';
  if (route === '/login' || route === '/register') return '0.5';
  return '0.7';
}

function routeChangeFreq(route) {
  if (route === '/') return 'weekly';
  if (route === '/login' || route === '/register') return 'monthly';
  return 'monthly';
}

function shouldIncludeRoute(route) {
  if (!route.startsWith('/')) return false;
  if (route.includes(':') || route.includes('*')) return false;
  if (EXCLUDED_ROUTES.has(route)) return false;
  return true;
}

function extractPublicRoutes(appSource) {
  const routes = new Set(['/']);
  const pathRegex = /\bpath\s*=\s*["']([^"']+)["']/g;

  for (const match of appSource.matchAll(pathRegex)) {
    const route = match[1].trim();
    if (shouldIncludeRoute(route)) {
      routes.add(route);
    }
  }

  return Array.from(routes).sort((a, b) => {
    if (a === '/') return -1;
    if (b === '/') return 1;
    return a.localeCompare(b);
  });
}

function generateSitemap(siteUrl, routes) {
  const lastmod = new Date().toISOString().slice(0, 10);
  const urls = routes
    .map((route) => {
      const loc = route === '/' ? `${siteUrl}/` : `${siteUrl}${route}`;
      return `  <url>
    <loc>${loc}</loc>
    <lastmod>${lastmod}</lastmod>
    <changefreq>${routeChangeFreq(route)}</changefreq>
    <priority>${routePriority(route)}</priority>
  </url>`;
    })
    .join('\n');

  return `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${urls}
</urlset>
`;
}

function main() {
  loadBuildEnv();
  const siteUrl = normalizeSiteUrl(process.env.SITE_URL || process.env.VITE_SITE_URL);
  const appSource = fs.readFileSync(appFile, 'utf-8');
  const routes = extractPublicRoutes(appSource);
  const sitemap = generateSitemap(siteUrl, routes);

  fs.writeFileSync(outputFile, sitemap, 'utf-8');
  // eslint-disable-next-line no-console
  console.log(`Generated sitemap with ${routes.length} route(s) at ${outputFile}`);
}

main();
