function arrayBufferToBase64(buf: ArrayBuffer): string {
  const bytes = new Uint8Array(buf);
  let str = "";
  const chunkSize = 0x8000;

  for(let i = 0; i < bytes.length; i+= chunkSize) {
    str += String.fromCharCode(...bytes.subarray(i, i + chunkSize));
  }

  return btoa(str);
}

function base64ToArrayBuffer(base64: string): ArrayBuffer {
  const binary = atob(base64);
  const len = binary.length;
  const bytes = new Uint8Array(len);

  for(let i = 0; i < len; i++){
    bytes[i] = binary.charCodeAt(i);
  }

  return bytes.buffer;
}

export async function generateKey(): Promise<CryptoKey> {
  return crypto.subtle.generateKey(
    {
      name: "AES-GCM",
      length: 256,
    },
    true,
    ["encrypt", "decrypt"]
  );
}

export async function exportKey(key: CryptoKey): Promise<string> {
  const raw = await crypto.subtle.exportKey("raw", key);
  return arrayBufferToBase64(raw);
}

export async function encryptFile(
  file: File,
  key: CryptoKey
): Promise<{
  encryptedBlob: Blob;
  iv: Uint8Array;
}> {
  const iv = crypto.getRandomValues(new Uint8Array(12));

  const fileBuffer = await file.arrayBuffer();

  const encryptedBuffer = await crypto.subtle.encrypt(
    {
      name: "AES-GCM",
      iv,
    },
    key,
    fileBuffer
  );

  return {
    encryptedBlob: new Blob([encryptedBuffer]),
    iv,
  };
}

export async function decryptFile(
  encrypted: ArrayBuffer,
  key: CryptoKey,
  iv: Uint8Array
): Promise<Blob> {
  const realIV = new Uint8Array(iv);
  const decryptedBuffer = await crypto.subtle.decrypt(
    {
      name: "AES-GCM",
      iv: realIV,
    },
    key,
    encrypted
  );

  return new Blob([decryptedBuffer]);
}

export async function importKeyFromBase64(base64: string): Promise<CryptoKey> {
  const raw = base64ToArrayBuffer(base64);
  return crypto.subtle.importKey("raw", raw, { name: "AES-GCM" }, true, [
    "encrypt",
    "decrypt",
  ]);
}

export async function generateRSAKeypair(): Promise<CryptoKeyPair> {
  return crypto.subtle.generateKey(
    {
      name: "RSA-OAEP",
      modulusLength: 2048,
      publicExponent: new Uint8Array([0x01, 0x00, 0x01]),
      hash: "SHA-256",
    },
    true,
    ["encrypt", "decrypt"]
  );
}

export async function exportPublicKeyToBase64(pub: CryptoKey): Promise<string> {
  const spki = await crypto.subtle.exportKey("spki", pub);
  return arrayBufferToBase64(spki);
}

export async function exportPrivateKeyToBase64(priv: CryptoKey): Promise<string> {
  const pkcs8 = await crypto.subtle.exportKey("pkcs8", priv);
  return arrayBufferToBase64(pkcs8);
}

export async function importPublicKeyFromBase64(base64: string): Promise<CryptoKey> {
  const ab = base64ToArrayBuffer(base64);
  return crypto.subtle.importKey(
    "spki",
    ab,
    { name: "RSA-OAEP", hash: "SHA-256" },
    true,
    ["encrypt"]
  );
}

export async function importPrivateKeyFromBase64(base64: string): Promise<CryptoKey> {
  const ab = base64ToArrayBuffer(base64);
  return crypto.subtle.importKey(
    "pkcs8",
    ab,
    { name: "RSA-OAEP", hash: "SHA-256" },
    true,
    ["decrypt"]
  );
}

export async function wrapAESKeytoBase64(aesKey: CryptoKey, recipientPublicKey: CryptoKey): Promise<string> {
  const wrapped = await crypto.subtle.wrapKey(
    "raw",
    aesKey,
    recipientPublicKey,
    { name: "RSA-OAEP" }
  );
  return arrayBufferToBase64(wrapped);
}

export async function unwrapAESKeyFromBase64(wrappedBase64: string, recipientPrivateKey: CryptoKey): Promise<CryptoKey> {
  const wrappedAb = base64ToArrayBuffer(wrappedBase64);
  const unwrapped = await crypto.subtle.unwrapKey(
    "raw",
    wrappedAb,
    recipientPrivateKey,
    { name: "RSA-OAEP" },
    { name: "AES-GCM", length: 256 },
    true,
    ["encrypt", "decrypt"]
  );
  return unwrapped;
}
