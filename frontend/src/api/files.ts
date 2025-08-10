import axios from "axios";

export async function uploadEncryptedFile(formData: FormData): Promise<void> {
    await axios.post("/api/files/upload", formData, {
        headers: { "Content-Type": "multipart/form-data" }
    });
}

export async function getFileMeta(fileId: string): Promise<{
    iv: string,
    wrappedKey: string;
    originalName: string;
}> {
    const res = await axios.get(`api/files/${fileId}/meta`);
    return res.data;
}

export async function getEncryptedFileData(fileId: string): Promise<ArrayBuffer> {
    const res = await axios.get(`api/files/${fileId}/download`, {
        responseType: "arraybuffer"
    });
    return res.data;
}