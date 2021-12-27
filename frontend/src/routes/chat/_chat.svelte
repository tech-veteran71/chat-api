<script lang="ts">
    import type { ScreenChat } from "./_store";
    import { readingByChatID, setChatAsRead } from "../api/_info";
    import { selectedChatID } from "./_store";

    export let chat: ScreenChat;

    $: label = chat.chatName || chat.senderName || chat.id;
    $: isSelected = ($selectedChatID === chat.id);
    $: isUnread = (chat.unread > 0);
    $: unreadCount = (chat.unread > 9) ? "+" : chat.unread;
    $: isReading = !!$readingByChatID[chat.id];

    $: bodyClass = isSelected && isUnread ? "bg-danger text-white"
                 : isSelected ? "bg-success text-white"
                 : isUnread ? "" : "bg-light";
    $: cardClass = isUnread ? "border-danger"
                 : isSelected ? "border-success"
                 : "border-muted text-muted";

    function onClick() {
        console.dir(chat);
        if (isSelected) {
            if (isUnread) {
                setChatAsRead(chat.id, chat.time);
            }
        } else {
            selectedChatID.set(chat.id);
        }
    }
</script>

<style>
    .spinner-border-badge {
        width: .75rem;
        height: .75rem;
    }
</style>

<div class="card {cardClass} border-3 mb-2">
    <div style="cursor: pointer" class="card-body py-1 {bodyClass} position-relative" on:click={onClick} title={JSON.stringify(chat, null, 2)}>
        {label}
        {#if isUnread}
            <span class="position-absolute top-0 start-100 translate-middle badge rounded-pill bg-danger">
                {#if isReading}
                    <div class="spinner-border spinner-border-badge" role="status">
                        <span class="visually-hidden">Marking chat as read</span>
                    </div>
                {:else}
                    {unreadCount}<span class="visually-hidden">unread messages</span>
                {/if}
            </span>
        {/if}
    </div>
</div>
