import type { ChatInfo } from "../api/_info"
import type { Message } from "$lib/whatsapp/chat-api/types"
import { derived, writable } from "svelte/store"
import { chatInfo, chatInfoByChatID } from "../api/_info"
import { messagesByChatID, messagesByID, message$ } from "../api/_messages"

// One screen contact.
export interface ScreenChat {
    id: string
    chatName: string
    senderName: string
    time: number // timestamp of last received message
    unread: number // number of unread messages
}

export const MessageIsRead = Symbol('MessageIsRead')

// One screen message.
export interface ScreenMessage extends Message {
    [MessageIsRead]?: boolean
}

// ID of the currently selected chat.
export const selectedChatID = writable('')

// Contact filter
export const screenChatsFilter = writable('')

// ScreenChat's - Left div.
export const screenChats = derived([message$, chatInfo, screenChatsFilter], ([$message, _chatInfo, screenChatsFilter]) => {
    console.log('Updating screenChats.')

    const chats = [] as ScreenChat[]

    for (const [chatID, messages] of Object.entries(messagesByChatID)) {
        const chat: ScreenChat = {
            id: chatID,
            chatName: '',
            senderName: '',
            time: 0,
            unread: 0,
        }
        // Get the time of the last read message.
        const info = chatInfoByChatID[chatID] ?? {} as ChatInfo
        const lastReadMessage = info.readBefore ?? 0
        // Count the number of unread messages.
        for (const message of Object.values(messages)) {
            if (message.fromMe) {
                continue // I read what I send
            }
            if (message.time > lastReadMessage) {
                chat.unread += 1
            }
        }
        // TODO: Messages are not ordered by timestamp yet here.
        // TODO: Use for starting at the last message instead of going through all messages.
        // Place the last chatName on Chat.
        for (const message of Object.values(messages)) {
            if (typeof message.chatName === 'string') {
                chat.chatName = message.chatName.trim()
            }
        }
        // Place the last senderName on Chat.
        for (const message of Object.values(messages)) {
            if (typeof message.senderName === 'string') {
                chat.senderName = message.senderName.trim()
            }
        }
        // Place the last time received on Chat.
        for (const message of Object.values(messages)) {
            if (message.fromMe || typeof message.time !== 'number') {
                continue
            }
            if (message.time > chat.time) {
                chat.time = message.time
            }
        }
        // Apply filter to the Chat list.
        if (!chat.chatName.includes(screenChatsFilter) && !chat.senderName.includes(screenChatsFilter)) {
            continue
        }
        // Show Chat to the user.
        chats.push(chat)
    }
    // TODO: Allow user to select the sorting algorithm.
    chats.sort((a, b) => b.time - a.time)

    return chats
})

// ScreenMessage's - Right div.
export const screenMessages = derived([message$, chatInfo, selectedChatID], ([$message, _chatInfo, selectedChatID]) => {
    console.log('Updating screenMessages %o.', $message.length)

    let messages: ScreenMessage[]

    // If a chat has been selected, show all messages in that chat.
    // Otherwise, show the most recent messages from all chats.
    if (selectedChatID) {
        const messagesByID = messagesByChatID[selectedChatID]
        if (!messagesByID) return []
        messages = Object.values(messagesByID)
        messages.sort((a, b) => b.time - a.time || b.messageNumber - a.messageNumber)
    } else {
        messages = [...$message]
        messages.sort((a, b) => b.time - a.time || b.messageNumber - a.messageNumber)
        messages = messages.slice(0, 250)
    }
    // Check whether each message has been read.
    for (const message of messages) {
        const chatID = message.chatId
        if (!chatID) {
            continue
        }
        const info = chatInfoByChatID[chatID]
        if (!info) {
            continue
        }
        const readBefore = info.readBefore
        if (!readBefore) {
            continue
        }
        if (message.time > readBefore) {
            continue
        }
        message[MessageIsRead] = true
    }

    return messages
})
