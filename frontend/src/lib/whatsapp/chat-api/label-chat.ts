import { getAuthorization } from "$lib/login";

export interface LabelChat {
    labelId: string
    chatId?: string
    phone?: string
}

export interface LabelChatResponse {
    chatId: string
    result: string
}

export function labelChat(options: LabelChat): Promise<LabelChatResponse> {
    return fetch(`/chat-api/labelChat`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `${getAuthorization()}`,
        },
        body: JSON.stringify(options)
    }).then(res => res.json())
}
