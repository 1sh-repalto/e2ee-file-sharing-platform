import axios from "axios";

export async function getPublicKeyBase64(recipientId: string): Promise<string> {
    const res = await axios.get(`api/users/${recipientId}/publicKey`);
    return res.data.publicKeyBase64;
}