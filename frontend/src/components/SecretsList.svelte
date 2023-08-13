<script lang="ts">
    import {Listgroup, ListgroupItem, P} from "flowbite-svelte";
    import type {Secret} from "$lib/secret";
    export let secrets: Secret[]
</script>

<Listgroup active class="w-full lg:w-1/2 2xl:w-2/5">
    {#each secrets as secret}
        <ListgroupItem class="">
            <P weight="semibold">{secret.key}</P>
            &nbsp;
            <P italic>[{secret.url}]</P>
        </ListgroupItem>
    {/each}
</Listgroup>