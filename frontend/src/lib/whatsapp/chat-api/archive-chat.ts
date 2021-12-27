import { getAuthorization } from "$lib/login";

export interface ArchiveChat {
    chatId?: string
    phone?: string
}

export interface ArchiveChatResponse {
    chatId: string
    result: string
}

export function archiveChat(options: ArchiveChat): Promise<ArchiveChatResponse> {
    return fetch(`/chat-api/archiveChat`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `${getAuthorization()}`,
        },
        body: JSON.stringify(options)
    }).then(res => res.json())
}
