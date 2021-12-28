import type { ChatInfo } from "../api/_info"
import type { Message } from "$lib/whatsapp/chat-api/types"
import { derived, writable } from "svelte/store"
import { chatInfo, chatInfoByChatID } from "../api/_info"
import { messagesByChatID, message$ } from "../api/_messages"
import { dialog$ } from "../api/_dialogs"

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

export const allScreenChats = derived([dialog$, message$, chatInfo, screenChatsFilter], ([$dialog, $message, _chatInfo, screenChatsFilter]) => {
    console.log('Updating screenChats %o + %o.', $dialog.length, $message.length)

    const chats = [] as ScreenChat[]
    const chatByID = {} as Record<string, ScreenChat>

    function getChat(chatID: string): ScreenChat {
        let chat = chatByID[chatID]
        if (chat === undefined) {
            chat = {
                id: chatID,
                chatName: '',
                senderName: '',
                time: 0,
                unread: 0,
            }
            chatByID[chatID] = chat
            chats.push(chat)
        }
        return chatByID[chatID]
    }

    // Add chats from messages.
    for (const [chatID, messages] of Object.entries(messagesByChatID)) {
        const chat = getChat(chatID)
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
    }

    // Add chats from dialogs.
    for (const dialog of $dialog) {
        const chat = getChat(dialog.id)
        chat.chatName = dialog.name
        if (chat.time === 0) {
            chat.time = dialog.last_time
        }
    }

    // TODO: Allow user to select the sorting algorithm.
    chats.sort((a, b) => b.time - a.time)

    return chats
})

// ScreenChat's - Left div.
export const screenChats = derived([allScreenChats, screenChatsFilter], ([chats, screenChatsFilter]) => {
    if (screenChatsFilter === '') {
        return chats
    }
    return chats.filter(filterChat)

    // Apply filter to the Chat list.
    function filterChat(chat: ScreenChat): boolean {
        return chat.chatName.includes(screenChatsFilter)
            || chat.senderName.includes(screenChatsFilter)
    }
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
