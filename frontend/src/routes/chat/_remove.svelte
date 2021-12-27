<script lang="ts">
    import type { ScreenMessage } from "./_store";
    import type { DeleteMessageResponse } from "$lib/whatsapp/chat-api";
    import { deleteMessage } from "$lib/whatsapp/chat-api";

    export let message: ScreenMessage;

    $: messageID = message.id;

    let error = '';
    let deleting = false;
    let deleted = false;

    function onDelete() {
        deleting = true;
        deleteMessage({ messageId: messageID })
            .then(onSuccess)
            .catch(handleError("Unable to remove message"))
            .finally(() => deleting = false);
            
        function onSuccess(response: DeleteMessageResponse) {
            console.dir(response);
            deleted = true;
        }
    }

    function handleError(text: string) {
        return (err: Error) => { error = `${text}:\n${err?.message ?? err}` };
    }
</script>

{#if deleted}
    <div class="alert alert-danger alert-dismissible mb-0">
        The message was removed from Chat-API!
        <button type="button" class="btn-close" aria-label="Close" on:click />
    </div>
{:else}
    <div class="card">
        <div class="card-body">
            <!-- Warning message -->
            <div class="mb-3">Are you sure you want to remove this message?<br /><br />The message will be removed from the Chat-API server, but it will be kept on this server.<br />If you try to remove more than one message at once, only the first may be removed from the other WhatsApp instance!</div>
            <!-- Buttons row -->
            <div class="row justify-content-between">
                <div class="col-auto">
                    <!-- Delete button -->
                    {#if deleting}<button class="btn btn-warning"><span class="spinner-border spinner-border-sm"></span> Removing</button>
                    {:else}<button class="btn btn-danger" on:click={onDelete}>Remove</button>{/if}
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
{/if}
