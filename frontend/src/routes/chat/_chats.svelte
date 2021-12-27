<script>
    import { screenChats, screenChatsFilter, selectedChatID } from "./_store";
    import Chat from "./_chat.svelte";

    $: isSelected = ($selectedChatID === "");
    $: bodyClass = isSelected ? "bg-secondary text-white" : "bg-light";
    $: cardClass = isSelected ? "border-secondary" : "border-muted text-muted";

    function onShowAll() {
        selectedChatID.set('');
    }
</script>

<input type="text" class="form-control my-3" placeholder="Chat filter" bind:value={$screenChatsFilter} />

<div class="card {cardClass} border-3 mb-2">
    <div on:click={onShowAll} style="cursor: pointer" class="card-body py-1 {bodyClass} position-relative">
        <i>Show All Messages</i>
    </div>
</div>

{#each $screenChats as chat (chat.id)}
    <Chat {chat} />
{/each}
