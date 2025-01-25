<script lang="ts">
    import Login from "./login.svelte";
    import Secrets from "./secrets.svelte";

    let isAuthenticated = $state(false);

    let { data } = $props();
    isAuthenticated = data.isAuthenticated;
</script>

{#if isAuthenticated }
    <Secrets bind:isAuthenticated/>
{:else}
    <Login />
{/if}
