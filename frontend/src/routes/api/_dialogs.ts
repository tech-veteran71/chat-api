import type { Dialog } from "$lib/whatsapp/chat-api/types"
import { getAuthorization } from "$lib/login"
import { writable } from "svelte/store"
import { syncStore } from "$lib/util"
import { f } from "$lib/fetch"

// All dialogs received from the server.
export const dialog$ = writable([] as Dialog[]) // [Dialog]
export const dialogsByID = {} as Record<string, Dialog> // {ID: [Dialog]}

let lastRowID: number = 0

function updateStores(newDialogs: Dialog[]): void {
    if (newDialogs.length === 0) {
        return
    }
    // Update row ID.
    for (const dialog of newDialogs) {
        if (dialog.__rowID > lastRowID) {
            lastRowID = dialog.__rowID
        }
    }
    // Update dialogsByID.
    for (const dialog of newDialogs) {
        dialogsByID[dialog.id] = dialog
    }
    // Update the store.
    dialog$.set(Object.values(dialogsByID))
}

export interface GetDialogsJSON {
    chats: Dialog[]
}

export async function getDialogs(): Promise<void> {
    console.log(`Getting dialogs.`)

    const params = new URLSearchParams()
    params.set('id', lastRowID.toString())
    params.set('wait', 'true')

    const headers = {}
    headers['Authorization'] = getAuthorization()

    const res = await f(`/api/chats/all?${params}`, { headers, timeout: 60000 })
    // TODO: Handle errors.

    const json: GetDialogsJSON = await res.json()
    updateStores(json.chats)
}

export const sync$ = syncStore('chats', getDialogs)
