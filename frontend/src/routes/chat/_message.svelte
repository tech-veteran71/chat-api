<script lang="ts">
    import type { ScreenMessage } from "./_store";
    import { MessageIsRead } from "./_store";
    import { selectedChatID } from "./_store";
    import moment from "moment";
    import { ackMap } from "./helper";
    import Reply from "./_reply.svelte";
    import Remove from "./_remove.svelte";

    export let message: ScreenMessage;

    const HIDDEN = 0;
    const REPLY = 1;
    const QUESTION = 2;
    const REMOVE =  3;
    let footer = HIDDEN;

    $: timestamp = moment.unix(message.time).format("DD/MM/YYYY HH:mm");
    $: color = message.fromMe ? "secondary"
             : message[MessageIsRead] ? "success"
             : "danger";
    $: ack = message.ack ?? '';
    $: ackLabel = ackMap[ack] ?? '?';

    function onClick() {
        console.dir(message);
    }

    function onHeader() {
        selectedChatID.set(message.chatId);
    }

    function onImage() {
        window.open(message.body, `wa-image-${message.id}`);
    }

    function onReply() {
        footer = REPLY;
    }

    function onQuestion() {
        if (footer === QUESTION) {
            location.hash = ('#' + message.quotedMsgId)
            footer = HIDDEN;
        } else {
            footer = QUESTION;
        }
    }

    function onRemove() {
        footer = REMOVE;
    }

    function onClose() {
        footer = HIDDEN;
    }
</script>

<style>
    img.fit {
        max-width: 100%;
        outline: solid 1px #ccc;
    }
</style>

<div class="card border-{color} mb-1" id={message.id} on:click={onClick}>
    <!-- Message header -->
    <div class="card-header text-white bg-{color}" style="cursor: pointer;" on:click={onHeader}>
        <div class="row">
            <!-- Sender or author -->
            <div class="col">
                {#if message.senderName}{message.senderName}
                {:else}<i>{message.author}</i>{/if}
            </div>
            <!-- Message status -->
            <div class="col-auto fw-bold" title={ack}>{ackLabel}</div>
            <!-- Timestamp -->
            <div class="col text-end">
                {timestamp}
            </div>
        </div>
    </div>
    <!-- Message body -->
    <div class="card-body" title={JSON.stringify(message, null, 2)}>
        {#if message.type === "chat"}
            <div style="white-space: pre-line">{message.body}</div>
        {:else if message.type === "vcard"}
            <div style="white-space: pre-line">{message.body}</div>
        {:else if message.type === "document"}
            {#if message.caption != null}<div class="mb-3" style="white-space: pre-line">{message.caption}</div>{/if}
            <a href="{message.body}" target="_blank">{message.body}</a>
        {:else if message.type === "audio" || message.type === "ptt"}
            <audio controls><source src={message.body} />Your browser does not support the audio tag.</audio>
        {:else if message.type === "image"}
            {#if message.caption != null}<div class="mb-3" style="white-space: pre-line">{message.caption}</div>{/if}
            <div><img class="fit" src={message.body} on:click={onImage} alt="Chat" /></div>
        {:else}
            {JSON.stringify(message, null, 2)}
        {/if}
        <!-- Line after message body -->
        <div class="row justify-content-end pt-3">
            <div class="col-auto">
                <!-- Question button -->
                {#if message.quotedMsgId}<button class="btn btn-outline-secondary" on:click={onQuestion}>Question</button>{/if}
                <!-- Remove button -->
                {#if footer !== REMOVE}<button class="btn btn-outline-danger" on:click={onRemove}>Remove</button>{/if}
                <!-- Reply button -->
                {#if footer !== REPLY}<button class="btn btn-outline-primary" on:click={onReply}>Reply</button>{/if}
            </div>
        </div>
    </div>
    <!-- Message commands -->
    {#if footer !== HIDDEN}
        <div class="card-footer">
            {#if footer === REPLY}
                <Reply message={message} on:click={onClose} />
            {:else if footer === REMOVE}
                <Remove message={message} on:click={onClose} />
            {:else}
                <div style="white-space: pre-line">{message.quotedMsgBody}</div>
            {/if}
        </div>
    {/if}
</div>
