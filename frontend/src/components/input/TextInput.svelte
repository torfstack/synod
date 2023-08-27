<script lang="ts">
    import {Helper, Input, Label} from "flowbite-svelte";

    export let value: string
    export let label: string
    export let handleKeyPress: ((event: KeyboardEvent) => void) | null = null
    export let error: boolean | null = null
    export let errorText: string | null = null
    export let required: boolean | null = null
    export let placeholder: string | null = null

    function clearError() {
        error = false
    }
</script>

<div>
    <Label color={ error ? "red" : "gray" } class="block mb-2">{label}</Label>
    <Input color={ error ? "red" : "base" } bind:value={value} type="text"
           on:keypress={handleKeyPress}
           on:focus={clearError}
           placeholder={placeholder}></Input>
    {#if required}
        <Helper color={ error ? "red" : "gray" } class="mt-2">
            { error && errorText ? errorText : "*required" }
        </Helper>
    {/if}
</div>
