import type { Message } from "./types"

export interface GetMessages {
    lastMessageNumber?: number
    last?: boolean
    chatId?: string
    limit?: number
    minTime?: number
    maxTime?: number
}

export interface GetMessagesResponse {
    lastMessageId: number
    messages: Message[]
}

export function getMessages(options: GetMessages): Promise<GetMessagesResponse> {
    const params = new URLSearchParams()
    if (options.lastMessageNumber != null) params.set('lastMessageNumber', options.lastMessageNumber.toString())
    if (options.last != null) params.set('last', 'true')
    if (options.chatId != null) params.set('chatId', options.chatId)
    if (options.limit != null) params.set('limit', options.limit.toString())
    if (options.minTime != null) params.set('min_time', options.minTime.toString())
    if (options.maxTime != null) params.set('max_time', options.maxTime.toString())
    return fetch(`/chat-api/messages?${params}`).then(res => res.json())
}
