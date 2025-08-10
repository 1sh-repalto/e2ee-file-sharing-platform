import { getEncryptedFileData, getFileMeta, uploadEncryptedFile } from "../api/files";
import { getPublicKeyBase64 } from "../api/users";
import {
  decryptFile,
  encryptFile,
  generateKey,
  importPrivateKeyFromBase64,
  importPublicKeyFromBase64,
  unwrapAESKeyFromBase64,
  wrapAESKeytoBase64,
} from "./crypto";

export async function uploadFage(
  selectedFile: File,
  recipientId: string
): Promise<void> {
  // 1️⃣ Generate AES key
  const aesKey = await generateKey();

  // 2️⃣ Encrypt file
  const { encryptedBlob, iv } = await encryptFile(selectedFile, aesKey);

  // 3️⃣ Get recipient's public RSA key
  const recipientPubKeyBase64 = await getPublicKeyBase64(recipientId);
  const recipientPubKey = await importPublicKeyFromBase64(
    recipientPubKeyBase64
  );

  // 4️⃣ Wrap AES key with RSA
  const wrappedKeyBase64 = await wrapAESKeytoBase64(aesKey, recipientPubKey);

  // 5️⃣ Send encrypted file + IV + wrapped key to backend
  const formData = new FormData();
  formData.append("file", encryptedBlob, selectedFile.name + ".enc");
  formData.append("iv", btoa(String.fromCharCode(...iv)));
  formData.append("wrappedKey", wrappedKeyBase64);

  await uploadEncryptedFile(formData);
}

export async function downloadFile(
  fileId: string,
  privateKeyBase64: string
): Promise<void> {
  const meta = await getFileMeta(fileId);

  const privateKey = await importPrivateKeyFromBase64(privateKeyBase64);

  const aesKey = await unwrapAESKeyFromBase64(meta.wrappedKey, privateKey);

  const encryptedData = await getEncryptedFileData(fileId);

  const iv = Uint8Array.from(atob(meta.iv), (c) => c.charCodeAt(0));
  const decryptedBlob = await decryptFile(encryptedData, aesKey, iv);

  const url = URL.createObjectURL(decryptedBlob);
  const a = document.createElement("a");
  a.href = url;
  a.download = meta.originalName || "downloaded_file";
  a.click();
  URL.revokeObjectURL(url);
}
