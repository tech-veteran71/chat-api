<script lang="ts">
    import type { SendMessage, CheckPhoneResponse } from "$lib/whatsapp/chat-api";
    import { readChat, sendMessage, checkPhone } from "$lib/whatsapp/chat-api";
    import { screenMessages, selectedChatID } from "./_store";
    import { getChatMessages } from "../api/_messages";
    import Message from "./_message.svelte";

    let phoneNumber = "";
    let messageBody = "";

    $: chatID = $selectedChatID
    $: canSend = messageBody !== ""
              && (chatID !== "" || phoneNumber !== "")

    let alertText = '';
    let alertColor = '';
    let reloading = false;
    let sending = false;
    let reading = false;
    let checking = false;

    function onReload() {
        reloading = true;
        getChatMessages(chatID, { update: true })
            .catch(handleError("Unable to update messages"))
            .finally(() => reloading = false);
    }

    function onSend() {
        const params: SendMessage = { body: messageBody };
        if (chatID) params.chatId = chatID;
        else params.phone = phoneNumber;
        sending = true;
        sendMessage(params)
            .catch(handleError("Unable to send message"))
            .finally(() => sending = false);
    }

    function onMarkAsRead() {
        if (!chatID) return;
        reading = true;
        readChat({ chatId: chatID })
            .catch(handleError("Unable to mark message as read"))
            .finally(() => reading = false);
    }

    function onCheck() {
        if (!phoneNumber) return;

        checking = true;
        checkPhone({ phone: phoneNumber })
            .then(handleResult)
            .catch(handleError("Unable to verify phone number"))
            .finally(() => checking = false);
        
        function handleResult(res: CheckPhoneResponse) {
            if (res.result === 'exists') {
                setInfo(`Phone number ${phoneNumber} is valid on WhatsApp.`);
            } else {
                const message = res.error ?? res.result ?? JSON.stringify(res);
                setError(`Phone number ${phoneNumber} was not found:\n${message}`);
            }
        }
    }

    function handleError(text: string) {
        return (err: Error) => { setError(`${text}:\n${err?.message ?? err}`) };
    }

    function setInfo(text: string) {
        alertText = text;
        alertColor = 'info';
    }

    function setError(text: string) {
        alertText = text;
        alertColor = 'danger';
    }

    function onAlertClose() {
        alertText = '';
    }
</script>

<div class="my-3">
    <!-- Phone number -->
    {#if chatID === ""}
        <div class="input-group mb-3">
            <input type="text" class="form-control" placeholder="Phone number" bind:value={phoneNumber}>
            <!-- Check button -->
            {#if checking}<button class="btn btn-warning"><span class="spinner-border spinner-border-sm"></span> Checking</button>
            {:else}<button class="btn btn-outline-secondary" on:click={onCheck} disabled={phoneNumber === ""}>Check</button>{/if}
        </div>
    {/if}
    <!-- Message area -->
    <textarea class="form-control mb-3" placeholder="New message" bind:value={messageBody} disabled={sending} />
    <div class="row justify-content-between mb-3">
        <div class="col-auto">
            <!-- Send button -->
            {#if sending}<button class="btn btn-warning"><span class="spinner-border spinner-border-sm"></span> Sending</button>
            {:else}<button class="btn btn-primary" on:click={onSend} disabled={!canSend}>Send</button>{/if}
        </div>
        <div class="col-auto">
            <!-- Read button -->
            {#if reading}<button class="btn btn-warning"><span class="spinner-border spinner-border-sm"></span> Reading</button>
            {:else if chatID !== ""}<button class="btn btn-outline-secondary" on:click={onMarkAsRead}>Read</button>{/if}
            {#if chatID !== ''}
                <!-- Reload button -->
                {#if reloading}<button class="btn btn-warning"><span class="spinner-border spinner-border-sm"></span></button>
                {:else}<button class="btn btn-outline-secondary" on:click={onReload}>🗘</button>{/if}
            {/if}
        </div>
    </div>
    <!-- Alert -->
    {#if alertText !== ""}<div class="alert alert-{alertColor} alert-dismissible" style="white-space: pre-line" role="alert">{alertText}
        <button type="button" class="btn-close" aria-label="Close" on:click={onAlertClose} /></div>{/if}
</div>

{#each $screenMessages as message (message.id)}
    <Message {message} />
{/each}
