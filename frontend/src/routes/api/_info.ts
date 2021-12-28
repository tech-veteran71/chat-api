import { browser } from "$app/env"
import { f } from "$lib/fetch"
import { getAuthorization } from "$lib/login"
import { get, writable } from "svelte/store"

export const chatInfo = writable([] as ChatInfo[])
export let chatInfoByChatID = {} as Record<string, ChatInfo>
export const readingByChatID = writable({} as Record<string, boolean>)

export interface ChatInfo {
    chatID: string
    readBefore: number
}

export interface GetInfoJSON {
    infos: ChatInfo[]
}

export async function getInfo(): Promise<void> {
    const headers = {}
    headers['Authorization'] = getAuthorization()
    const res = await fetch('/api/chat/info', { headers })
    const json = await res.json() as GetInfoJSON

    chatInfoByChatID = {}
    for (const info of json.infos) {
        chatInfoByChatID[info.chatID] = info
    }

    chatInfo.set(json.infos)
}

export async function setChatAsRead(chatID: string, timestamp: number): Promise<void> {
    const $readingByChatID = get(readingByChatID)
    if ($readingByChatID[chatID]) return

    console.log("setChatAsRead(%o, %o)", chatID, timestamp)

    try {
        $readingByChatID[chatID] = true
        readingByChatID.set($readingByChatID)

        const headers = {}
        headers['Authorization'] = getAuthorization()
        const res = await f(`/api/chat/info/read`, {
            method: 'POST', headers, timeout: 10000,
            body: JSON.stringify({ chatID, timestamp }),
        })
        if (!res.ok) {
            throw new Error(`${res.status} ${res.statusText}`)
        }
    } catch (err) {
        console.log("Could not set chat as read.")
        console.error(err)
        setTimeout(() => setChatAsRead.apply(this, arguments), 1000)
        return
    } finally {
        delete $readingByChatID[chatID]
        readingByChatID.set($readingByChatID)
    }

    let info = chatInfoByChatID[chatID]
    if (info == null) {
        info = { chatID, readBefore: timestamp }
        chatInfoByChatID[chatID] = info
        chatInfo.update(infos => [...infos, info])
    } else {
        info.readBefore = timestamp
        chatInfo.update(infos => infos)
    }
}

if (browser) {
    getInfo().catch(console.error)
}
