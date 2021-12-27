import { getAuthorization } from "$lib/login"
import type { Label } from "./types"

export interface GetLabelsResponse {
    labels: Label[]
}

export function getLabels(): Promise<GetLabelsResponse> {
    const headers = {}
    headers['Authorization'] = getAuthorization()
    return fetch(`/chat-api/labelsList`, { headers }).then(res => res.json())
}
