<script lang="ts">
    import {Listgroup, ListgroupItem, P} from "flowbite-svelte";
    import type {Secret} from "$lib/secret";
    interface Props {
        secrets: Secret[];
    }

    let { secrets }: Props = $props();
</script>

{#if secrets.length > 0}
    <Listgroup active class="w-full lg:w-1/2 2xl:w-2/5">
        {#each secrets as secret}
            <ListgroupItem class="">
                <P weight="semibold">{secret.key}</P>
                &nbsp;
                <P italic>[{secret.url}]</P>
            </ListgroupItem>
        {/each}
    </Listgroup>
{/if}
