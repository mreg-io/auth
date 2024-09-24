export function parseCookie() {
  return new Map(
    document.cookie.split("; ").map((v) => {
      const [key, ...val] = v.split(/=(.*)/s);
      return [decodeURIComponent(key), decodeURIComponent(val.join("="))];
    })
  );
}
