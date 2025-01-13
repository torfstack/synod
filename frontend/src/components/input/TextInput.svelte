<script lang="ts">
    import {Helper, Input, Label} from "flowbite-svelte";

    interface Props {
        value: string;
        label: string;
        handleKeyPress?: ((event: KeyboardEvent) => void) | null;
        error?: boolean | null;
        errorText?: string | null;
        required?: boolean | null;
        placeholder?: string | null;
    }

    let {
        value = $bindable(),
        label,
        handleKeyPress = null,
        error = $bindable(null),
        errorText = null,
        required = null,
        placeholder = null
    }: Props = $props();

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
