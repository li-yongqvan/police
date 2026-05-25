/**
 * Parallel page data loading (SvelteKit-style load, minimal surface).
 */
export async function loadPage(loaders) {
  const keys = Object.keys(loaders)
  const tasks = keys.map((key) => loaders[key]())
  const results = await Promise.all(tasks)
  return Object.fromEntries(keys.map((key, i) => [key, results[i]]))
}
