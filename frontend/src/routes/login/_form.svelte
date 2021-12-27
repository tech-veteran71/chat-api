<script lang="ts">
    export let data = {} as any;

    let username = data.username ?? "";
    $: data.username = username;

    let password = data.password ?? "";
    $: data.password = password;

    let remember = data.remember ?? false;
    $: data.remember = remember;
</script>

<div class="mb-3">
    <label class="form-label" for="login-username">User</label>
    <input type="text" class="form-control" name="username" id="login-username" bind:value={username} />
</div>

<div class="mb-3">
    <label class="form-label" for="login-password">Password</label>
    <input type="text" class="form-control" name="password" id="login-password" bind:value={password} />
</div>

<div class="form-check">
    <input class="form-check-input" type="checkbox" id="login-remember" bind:checked={remember} />
    <label class="form-check-label" for="login-remember">Remember me.</label>
</div>
