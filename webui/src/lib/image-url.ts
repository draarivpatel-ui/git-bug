// The spelling in values.md §9, restricted to the HTTPS port allowed by §9.1.
// Check the input before the browser normalizes it. The page's img-src policy
// enforces the scheme and port again when the browser follows a redirect.
const imageURL =
  /^https:\/\/(?:[a-zA-Z0-9-]{1,63}\.)+[a-zA-Z][a-zA-Z0-9-]{0,62}(?::0*443)?(?:[/?#](?:[a-zA-Z0-9\-._~:/?#[\]@!$&'()*+,;=]|%[a-fA-F0-9]{2})*)?$/;

export function remoteImageURL(src: string | undefined): string | undefined {
  return src && src.length <= 2048 && imageURL.test(src) ? src : undefined;
}

// Keep simple srcset candidates and their optional width or density descriptors.
// Rebuild the list so that the browser sees the same candidate boundaries we checked.
export function remoteImageSet(srcSet: string | undefined): string | undefined {
  if (!srcSet) return undefined;
  const candidates = [];
  for (const candidate of srcSet.split(",")) {
    const [url, descriptor, ...extra] = candidate.trim().split(/\s+/);
    if (!remoteImageURL(url) || extra.length) return undefined;
    if (descriptor && !/^(?:\d+w|(?:\d+(?:\.\d+)?|\.\d+)x)$/.test(descriptor)) return undefined;
    candidates.push(descriptor ? `${url} ${descriptor}` : url);
  }
  return candidates.join(", ");
}
