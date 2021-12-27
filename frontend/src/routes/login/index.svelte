<script lang="ts">
    import { goto } from "$app/navigation";
    import { login, saving } from "$lib/login";
    import Form from "./_form.svelte";

    let data = {
        username: "",
        password: "",
        remember: false,
    };

    let alertText = "";

    async function onLogin() {
        try {
            await login(data.username, data.password, data.remember);
            goto("/chat");
        } catch (error) {
            alertText = `${error.message ?? error}`;
        }
    }
</script>

<div class="container pt-5">
    {#if $saving}
        <div class="text-center mt-5">
            <div class="alert alert-info d-inline">
                <span class="spinner-border spinner-border-sm"></span>
                <span>Logging in…</span>
            </div>
        </div>
    {:else}
        <div class="row justify-content-center">
            <div class="col-md-6 col-xl-4">
                <div class="card">
                    <div class="card-body">
                        <!-- Login form -->
                        <Form bind:data />
                        <!-- Alert -->
                        {#if alertText}<div class="alert alert-danger mt-3 mb-0" role="alert">{alertText}</div>{/if}
                    </div>
                    <div class="card-footer text-center">
                        <!-- Login button -->
                        <button class="btn btn-primary" on:click={onLogin}>Login</button>
                    </div>
                </div>
            </div>
        </div>
    {/if}
</div>