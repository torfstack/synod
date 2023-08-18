<script lang="ts">
    import {Banner, Heading, Mark, P, Span, DarkMode} from "flowbite-svelte";
    export let text: string;
</script>

<div>
    <Banner dismissable={false}>
        <div class="flex flex-row justify-center">
            <div class="">
                <Heading class="text-center mb-2 text-5xl" tag="h1"><Span class="text-primary-600 dark:text-primary-600 font-mono">Kay</Span>Vault</Heading>
                <P class="text-center">{text} <Mark>secrets</Mark></P>
            </div>
            <div class="absolute right-4">
                <DarkMode />
            </div>
        </div>
    </Banner>
</div>
