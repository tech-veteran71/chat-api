import { getAuthorization } from "$lib/login";

export interface DeleteMessage {
    messageId?: string
}

export interface DeleteMessageResponse {
    sent: boolean
    message: string // "message deleted"
}

export function deleteMessage(options: DeleteMessage): Promise<DeleteMessageResponse> {
    return fetch(`/chat-api/deleteMessage`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `${getAuthorization()}`,
        },
        body: JSON.stringify(options)
    }).then(res => res.json())
}
