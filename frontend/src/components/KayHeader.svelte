<script lang="ts">
    import {Banner, Heading, Mark, P, Span} from "flowbite-svelte";
    export let text: string;
</script>

<div>
    <Banner dismissable={false}>
        <div>
            <Heading class="text-center mb-2 text-5xl" tag="h1"><Span class="text-primary-600 font-mono">Kay</Span>Vault</Heading>
            <P class="text-center">{text} <Mark>secrets</Mark></P>
        </div>
    </Banner>
</div>
