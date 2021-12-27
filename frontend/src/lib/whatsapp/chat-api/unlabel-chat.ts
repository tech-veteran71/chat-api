import { getAuthorization } from "$lib/login";

export interface UnlabelChat {
    labelId: string
    chatId?: string
    phone?: string
}

export interface UnlabelChatResponse {
    chatId: string
    result: string
}

export function unlabelChat(options: UnlabelChat): Promise<UnlabelChatResponse> {
    return fetch(`/chat-api/unlabelChat`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `${getAuthorization()}`,
        },
        body: JSON.stringify(options)
    }).then(res => res.json())
}
