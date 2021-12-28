import { goto } from "$app/navigation"

export const ErrTimeout = Error("Network or server timeout")

export interface FOptions extends RequestInit {
    timeout: number
}

// f is a wrapper around fetch that adds a timeout option.
export function f(url, options = {} as FOptions) {
    const controller = new AbortController()

    if (options.signal) {
        options.signal.addEventListener('abort', () => {
            console.log("Provided signal was aborted.")
            controller.abort()
        })
    }

    const newOptions = {
        ...options,
        signal: controller.signal,
    }

    const timeout = setTimeout(() => {
        console.log("Request timed out.")
        controller.abort()
    }, options.timeout)

    return fetch(url, newOptions)
        .then(handleLogin)
        .catch(handleAbort)
        .finally(() => clearTimeout(timeout))
}

function handleLogin(res: Response): Response {
    if (res.status === 401) {
        console.log(`Got ${res.status} ${res.statusText} from the server.`)
        console.log(`Redirecting to login page.`)
        setTimeout(() => goto("/login"), 100)
    } else {
        return res
    }
}

function handleAbort(err: DOMException): never {
    if (err.name === 'AbortError') {
        throw ErrTimeout
    } else {
        throw err
    }
}
