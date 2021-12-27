import { getAuthorization } from "$lib/login"
import type { Dialog } from "./types"

export interface GetDialogs {
    limit?: number
    page?: number
    order?: string
}

export interface GetDialogsResponse {
    dialogs: Dialog[]
}

export function getDialogs(options: GetDialogs): Promise<GetDialogsResponse> {
    const params = new URLSearchParams()
    if (options.limit != null) params.set('limit', options.limit.toString())
    if (options.page != null) params.set('page', options.page.toString())
    if (options.order != null) params.set('order', options.order)
    const headers = {}
    headers['Authorization'] = getAuthorization()
    return fetch(`/chat-api/dialogs?${params}`, { headers }).then(res => res.json())
}
