<script lang="ts">
    import type { ScreenMessage } from "./_store";
    import { sendMessage } from "$lib/whatsapp/chat-api";

    export let message: ScreenMessage;

    $: chatID = message.chatId;
    $: messageID = message.id;

    let messageBody = "";

    let error = '';
    let sending = false;

    function onSend() {
        sending = true;
        sendMessage({body: messageBody, chatId: chatID, quotedMsgId: messageID})
            .catch(handleError("Unable to send message"))
            .finally(() => sending = false);
    }

    function handleError(text: string) {
        return (err: Error) => { error = `${text}:\n${err?.message ?? err}` };
    }
</script>

<div class="card">
    <div class="card-body">
        <!-- Message area -->
        <textarea class="form-control mb-3" placeholder="New reply" bind:value={messageBody} disabled={sending} />
        <!-- Buttons row -->
        <div class="row justify-content-between">
            <div class="col-auto">
                <!-- Send button -->
                {#if sending}<button class="btn btn-warning"><span class="spinner-border spinner-border-sm"></span> Sending</button>
                {:else}<button class="btn btn-primary" on:click={onSend} disabled={messageBody === ""}>Send</button>{/if}
            </div>
            <div class="col-auto">
                <!-- Close button -->
                <button class="btn btn-outline-danger" on:click>Close</button>
            </div>
        </div>
        <!-- Error alert -->
        {#if error !== ""}<div class="alert alert-danger mt-3 mb-0" style="white-space: pre-line" role="alert">{error}</div>{/if}
    </div>
</div>
