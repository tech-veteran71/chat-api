import { f } from "$lib/fetch"
import { getAuthorization } from "$lib/login"
import type { Message } from "$lib/whatsapp/chat-api/types"
import { writable } from "svelte/store"

// All messages received from the server.
export const message$ = writable([] as Message[]) // [Message]
export const messagesByID = {} as Record<string, Message> // {ID: [Message]}
export const messagesByChatID = {} as Record<string, Record<string, Message>> // {ChatID: {ID: [Message]}}
export const newestMessageTimeByChatID = {} as Record<string, number> // {ChatID: Timestamp}

// The ID of the last message fetched from the server.
// Function getMessages uses this to only fetch new messages.
let lastRowID: number = 0

// All messages received from the server should be given to this function.
// It updates all variables above and then updates the store.
function updateStores(newMessages: Message[]): void {
    if (newMessages.length === 0) {
        return
    }
    // Update lastMessageNumber.
    for (const message of newMessages) {
        if (message.__rowID > lastRowID) {
            lastRowID = message.__rowID
        }
    }
    // Update newestMessageTimeByChatID.
    for (const message of newMessages) {
        const time = newestMessageTimeByChatID[message.chatId] ?? 0
        if (message.time > time) {
            newestMessageTimeByChatID[message.chatId] = message.time
        }
    }
    // Update messagesByID.
    for (const message of newMessages) {
        messagesByID[message.id] = message
    }
    // Update messagesByChatID.
    for (const message of newMessages) {
        const chatID = message.chatId
        if (!chatID) continue
        (messagesByChatID[chatID] ??= {})[message.id] = message
    }
    // Update the store.
    message$.set(Object.values(messagesByID))
}

// GetMessagesJSON is the shape of our API response.
export interface GetMessagesJSON {
    messages: Message[]
}

export async function getMessages(): Promise<void> {
    console.log(`Getting messages after row ID %o.`, lastRowID)

    const params = new URLSearchParams()
    params.set('wait', 'true')
    params.set('id', lastRowID.toString())

    const headers = {}
    headers['Authorization'] = getAuthorization()

    const res = await f(`/api/messages/all?${params}`, { headers, timeout: 30000 })
    // TODO: Handle errors.

    const json: GetMessagesJSON = await res.json()
    updateStores(json.messages)
}

export async function getChatMessages(chatID: string): Promise<void> {
    console.log(`Getting messages from chat ID %o.`, chatID)

    const params = new URLSearchParams()
    params.set('chat_id', chatID)

    const headers = {}
    headers['Authorization'] = getAuthorization()

    const res = await f(`/api/messages/chat_id?${params}`, { headers, timeout: 30000 })
    // TODO: Handle errors.

    const json: GetMessagesJSON = await res.json()
    updateStores(json.messages)
}
