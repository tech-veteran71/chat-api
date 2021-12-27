import { getAuthorization } from "$lib/login";

export interface UnreadChat {
    chatId?: string
    phone?: string
}

export interface UnreadChatResponse {
    chatId: string
    result: string // "success"
}

export function unreadChat(options: UnreadChat): Promise<UnreadChatResponse> {
    return fetch(`/chat-api/unreadChat`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `${getAuthorization()}`,
        },
        body: JSON.stringify(options)
    }).then(res => res.json())
}
