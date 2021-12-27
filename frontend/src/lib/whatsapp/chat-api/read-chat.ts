import { getAuthorization } from "$lib/login";

export interface ReadChat {
    chatId?: string
    phone?: string
}

export interface ReadChatResponse {
    read: boolean
    chatId: string
    message: string
}

export function readChat(options: ReadChat): Promise<ReadChatResponse> {
    return fetch(`/chat-api/readChat`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `${getAuthorization()}`,
        },
        body: JSON.stringify(options)
    }).then(res => res.json())
}
