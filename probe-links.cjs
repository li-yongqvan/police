const { chromium } = require("playwright");
(async () => {
  const b = await chromium.launch({ headless: true });
  const p = await b.newPage({ viewport: { width: 390, height: 844 } });
  await p.goto("http://122.51.233.225:8080/latest", { waitUntil: "load", timeout: 45000 });
  await p.waitForTimeout(4000);
  const links = await p.$$eval("a[href*='/t/'], a.title", els => els.slice(0, 10).map(a => ({ href: a.href, cls: (a.className||"").toString().slice(0,60), text: (a.innerText||"").trim().slice(0,30) })));
  console.log(JSON.stringify(links, null, 1));
  await b.close();
})();