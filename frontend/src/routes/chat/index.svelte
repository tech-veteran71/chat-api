<script lang="ts">
    import { sync$ as messageSync } from "../api/_messages";
    import { sync$ as chatSync } from "../api/_dialogs";
    import { selectedChatID } from "./_store";
    import Messages from "./_messages.svelte";
    import Chats from "./_chats.svelte";

    $: console.dir($chatSync);
    $: console.dir($messageSync);

    // Scroll to the top when another chat is selected.
    let messagesDiv = {} as Element;
    $: $selectedChatID, messagesDiv.scrollTop = 0;
</script>

<div class="container-fluid">
    <div class="row vh-100">
        <div class="col-3 h-100 overflow-auto">
            <Chats />
        </div>
        <div bind:this={messagesDiv} class="col-9 h-100 overflow-auto">
            <Messages />
        </div>
    </div>
</div>
