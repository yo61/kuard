// Smoke test for the *built* bundle, deliberately not for the source.
//
// The React 19 migration shipped a bundle that compiled cleanly, served a 200, and rendered
// nothing: Babel emitted jsxDEV(), which React's production jsx-dev-runtime does not export, so
// the app threw during mount with an empty console. A test that imported ./app and rendered it
// would not have caught that -- vitest transforms source with its own esbuild pipeline, not with
// our webpack + Babel production config. The bug only exists in the emitted output, so that is
// what this loads.
//
// Requires `npm run build` first; the npm test script chains them.

import { readFileSync } from 'fs';
import { fileURLToPath } from 'url';
import { JSDOM } from 'jsdom';
import { beforeAll, describe, expect, it } from 'vitest';

const BUNDLE = fileURLToPath(new URL('../sitedata/built/bundle.js', import.meta.url));

// Mirrors the pageContext the Go template writes into the page. The root route (Request Details)
// is deliberately the one under test: it is the only page that issues no fetch on mount, so this
// needs no network stubbing.
const PAGE_CONTEXT = {
  urlBase: '',
  hostname: 'test-host',
  addrs: ['10.0.0.1'],
  version: 'blue',
  versionColor: '#000000',
  requestDump: 'GET / HTTP/1.1',
  requestProto: 'HTTP/1.1',
  requestAddr: '10.0.0.2:1234',
};

// The URL has to sit inside urlBase, because that is what kuard does: it only reports
// urlBase "/b" for a request that arrived at /b/. React Router renders nothing when the current
// location falls outside its basename, so a mismatch here tests nothing.
function mount(pageContext = PAGE_CONTEXT) {
  const dom = new JSDOM('<!doctype html><html><body><div id="root"></div></body></html>', {
    url: `http://localhost${pageContext.urlBase}/`,
    runScripts: 'outside-only',
    pretendToBeVisual: true,
  });

  dom.window.pageContext = pageContext;
  dom.window.eval(readFileSync(BUNDLE, 'utf8'));

  return dom;
}

// React 19 schedules work rather than rendering synchronously, and under jsdom that takes more
// than one macrotask tick. Poll to a deadline instead of guessing a delay: fast when it works,
// and it fails with "never mounted" rather than a flake when it doesn't.
async function waitForMount(dom, timeoutMs = 5000) {
  const root = dom.window.document.getElementById('root');
  const deadline = Date.now() + timeoutMs;

  while (Date.now() < deadline) {
    if (root.children.length > 0) return root;
    await new Promise((resolve) => setTimeout(resolve, 10));
  }

  throw new Error(`React never mounted into #root within ${timeoutMs}ms`);
}

describe('built bundle', () => {
  let bundle;

  beforeAll(() => {
    bundle = readFileSync(BUNDLE, 'utf8');
  });

  it('mounts React into #root', async () => {
    const dom = mount();
    const root = await waitForMount(dom);

    expect(root.children.length).toBeGreaterThan(0);
  });

  it('renders the page content the server handed it', async () => {
    const dom = mount();
    const root = await waitForMount(dom);

    const text = root.textContent;
    expect(text).toContain('test-host');
    expect(text).toContain('blue');
    expect(text).toContain('10.0.0.1');
  });

  it('renders the nav and routes to Request Details at /', async () => {
    const dom = mount();
    await waitForMount(dom);

    const doc = dom.window.document;
    const navLabels = [...doc.querySelectorAll('.nav-item')].map((el) => el.textContent);
    expect(navLabels).toContain('Request Details');
    expect(navLabels).toContain('MemQ Server');

    // Request Details is the "/" route, so its request dump should be on screen.
    expect(doc.getElementById('root').textContent).toContain('GET / HTTP/1.1');
  });

  it('marks only the active nav item', async () => {
    const dom = mount();
    await waitForMount(dom);

    const active = [...dom.window.document.querySelectorAll('.nav-item.active')];
    expect(active.map((el) => el.textContent)).toEqual(['Request Details']);
  });

  // urlBase is passed to React Router as its basename. kuard serves the app under "", /a, /b and
  // /c, and a broken basename silently breaks every link.
  it('honours a urlBase prefix in nav links', async () => {
    const dom = mount({ ...PAGE_CONTEXT, urlBase: '/b' });
    await waitForMount(dom);

    const hrefs = [...dom.window.document.querySelectorAll('.nav-item')].map((el) =>
      el.getAttribute('href'),
    );
    // React Router collapses basename + "/" to just the basename.
    expect(hrefs).toContain('/b');
    expect(hrefs).toContain('/b/-/memq');
    expect(hrefs.every((href) => href.startsWith('/b'))).toBe(true);
  });

  it('uses unprefixed links when urlBase is empty', async () => {
    const dom = mount();
    await waitForMount(dom);

    const hrefs = [...dom.window.document.querySelectorAll('.nav-item')].map((el) =>
      el.getAttribute('href'),
    );
    expect(hrefs).toContain('/');
    expect(hrefs).toContain('/-/memq');
  });

  // Guards the specific failure above: a production bundle must not reference the dev runtime.
  it('does not reference the development JSX runtime', () => {
    expect(bundle).not.toContain('jsxDEV');
  });
});
