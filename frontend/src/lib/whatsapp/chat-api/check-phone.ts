import { getAuthorization } from "$lib/login"

export interface CheckPhone {
    phone: string
}

export interface CheckPhoneResponse {
    result: string // "exists" | "Status of the account not equals authenticated"
    error: string // not in the API docs, but does happen
}

export function checkPhone(options: CheckPhone): Promise<CheckPhoneResponse> {
    const params = new URLSearchParams({ phone: options.phone })
    const headers = {}
    headers['Authorization'] = getAuthorization()
    return fetch(`/chat-api/checkPhone?${params}`, { headers }).then(res => res.json())
}
