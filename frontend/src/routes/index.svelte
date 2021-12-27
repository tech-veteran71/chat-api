<script>
    import { goto } from "$app/navigation"
    import { getToken } from "$lib/login"

    if (getToken()) goto("/chat")
    else goto("/login")
</script>
