<script lang="ts">
    import {Heading, Mark, P, Span} from "flowbite-svelte";
    export let text: string;
</script>

<Heading class="text-center mb-2 text-7xl" tag="h1"><Span class="font-mono" gradient>Kay</Span>Vault</Heading>
<P class="text-center">{text} <Mark>secrets</Mark></P>
