import { get, writable } from "svelte/store"
import { browser } from "$app/env"
import { f } from "./fetch"
import { sleep } from "./util"

// token is a store that keeps the current authentication token. It should be included in all API requests.
// It's an empty string when the user is not logged in yet.
export let token = writable(getStoredToken())

// saving is a boolean store that indicates whether the user is logging in.
export let saving = writable(false)

// getToken returns the current token.
export function getToken(): string {
    return get(token)
}

// getAuthorization returns the value of the Authorization header.
export function getAuthorization(): string {
    return `Bearer ${getToken()}`
}

// getStoredToken returns the token stored in sessionStorage or localStorage.
// Returns null when not running in the browser.
function getStoredToken(): string {
    if (!browser) return
    return sessionStorage.getItem('token') ?? localStorage.getItem('token') ?? ''
}

// login logs the user in.
export async function login(username: string, password: string, remember: boolean): Promise<void> {
    const params = {
        method: 'POST', timeout: 10000,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    }

    saving.set(true)
    try {
        const res = await f(`/api/login`, params)
        if (!res.ok) {
            throw await res.text()
        }
        const json = await res.json()
        if (typeof json.error === 'string') {
            throw json.error
        }
        if (typeof json.token !== 'string') {
            throw `Invalid server response`
        }
        console.log('Logged in with token %o', json.token)
        token.set(json.token)
        if (remember) {
            localStorage.setItem('token', json.token)
        } else {
            sessionStorage.setItem('token', json.token)
        }
    } catch (err) {
        console.dir(err)
        throw `${err}`
    } finally {
        saving.set(false)
    }
}

// logout logs the user out, by removing the authentication token from memory and storage.
// Function f will be called until 401 is received, at which point it will redirect to /login.
export async function logout(): Promise<void> {
    const params = {
        method: 'GET', timeout: 10000,
        headers: { Authorization: getAuthorization() }
    }

    token.set('')
    sessionStorage.removeItem('token')
    localStorage.removeItem('token')

    for (; ;) {
        try {
            const res = await f(`/api/logout`, params)
            if (res.status === 401) {
                return
            }
        } catch (err) {
            console.log("Could not log out.")
            console.error(err)
        } finally {
            await sleep(1000)
        }
    }
}
