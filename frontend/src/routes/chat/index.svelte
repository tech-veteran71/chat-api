<script lang="ts">
    import { getMessages } from "../api/_messages";
    import { ErrTimeout } from "$lib/fetch";
    import { onMount } from "svelte";
    import { selectedChatID } from "./_store";
    import Chats from "./_chats.svelte";
    import Messages from "./_messages.svelte";

    // Keep the messages store updated.
    let timer: any;
    onMount(() => {
        timer = setTimeout(onTimer, 0);
        return () => clearTimeout(timer);
    });
    async function onTimer() {
        let delay = 100; // ms
        try {
            // TODO: Tell the user that something went wrong.
            await getMessages()
        } catch (err) {
            if (err !== ErrTimeout) {
                console.log("Could not get new messages.")
                console.error(err)
                delay = 5000; // ms
            }
        } finally {
            timer = setTimeout(onTimer, delay);
        }
    }

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
