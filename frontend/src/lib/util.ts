import type { Readable } from "svelte/store"
import { readable } from "svelte/store"
import { ErrTimeout } from "./fetch"

export function sleep(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms))
}

export function syncStore(name: string, onRun: () => Promise<unknown>): Readable<string> {
    let timeout: NodeJS.Timeout

    return readable('', set => {
        console.log(`Starting to sync ${name}.`)

        run()

        function run() {
            onRun().then(onSuccess, onError)
        }

        function onSuccess() {
            timeout = setTimeout(run, 100)
            set('')
        }

        function onError(err: Error) {
            if (err === ErrTimeout) {
                onSuccess()
            } else {
                timeout = setTimeout(run, 60000)
                set(`Could not update ${name}: ${err}`)
            }
        }

        return () => {
            clearTimeout(timeout)
            console.log(`Stopped to sync ${name}.`)

        }
    })
}
