import type { Dialog, Message } from "$lib/whatsapp/chat-api/types"
import { writable } from "svelte/store"

export const dialogs = writable({} as Record<string, Dialog>)

export async function getDialogs() {
    console.log(`Getting dialogs.`)
    const res = await fetch(`/api/dialogs.json`)
    const json = await res.json()
    dialogs.update(msgs => ({
        ...msgs, ...json.dialogs.reduce((acc, msg) => {
            acc[msg.id] = msg
            return acc
        }, {})
    }))
}
