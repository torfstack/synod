<script lang="ts">
    import {Helper, Input, Label} from "flowbite-svelte";

    interface Props {
        value: string;
        label: string;
        handleKeyPress?: ((event: KeyboardEvent) => void);
        error?: boolean;
        errorText?: string;
        required?: boolean;
        placeholder?: string;
    }

    let {
        value = $bindable(),
        label,
        handleKeyPress,
        error = $bindable(),
        errorText,
        required,
        placeholder
    }: Props = $props();

    function clearError() {
        error = false
    }
</script>

<div>
    <Label class="block mb-2" color={ error ? "red" : "gray" }>{label}</Label>
    {#if handleKeyPress === undefined}
        <Input bind:value={value} color={ error ? "red" : "base" } on:focus={clearError}
               placeholder={placeholder}
               type="text"></Input>
    {:else}
        <Input bind:value={value} color={ error ? "red" : "base" } on:focus={clearError}
               on:keypress={handleKeyPress}
               placeholder={placeholder}
               type="text"></Input>
    {/if}
    {#if required}
        <Helper color={ error ? "red" : "gray" } class="mt-2">
            { error && errorText ? errorText : "*required" }
        </Helper>
    {/if}
</div>
