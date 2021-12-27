import { getAuthorization } from "$lib/login";

export interface SendMessage {
    body: string
    quotedMsgId?: string
    chatId?: string
    phone?: string
    mentionedPhones?: string[]
}

export interface SendMessageResponse {
    sent: boolean
    message: string // "Sent to 5512345678@c.us"
    id: string // "true_5512345678@c.us_3EB0C7F9412345678908"
    queueNumber: number // 1
}

export function sendMessage(options: SendMessage): Promise<SendMessageResponse> {
    return fetch(`/chat-api/sendMessage`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `${getAuthorization()}`,
        },
        body: JSON.stringify(options)
    }).then(res => res.json())
}
