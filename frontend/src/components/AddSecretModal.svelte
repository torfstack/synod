<script lang="ts">
    import {Secret} from "$lib/secret"
    import {Badge, Button, Input, Label, Modal, P} from "flowbite-svelte";
    import {slide} from "svelte/transition";

    export let uploadSecret: (n: Secret) => Promise<void>;
    export let openModal: boolean;

    let inputKey: string, inputValue: string, inputUrl: string, inputTag: string, inputTags: string[] = []

    function handleSecret() {
        let secret = new Secret(
            inputKey,
            inputValue,
            inputUrl,
            inputTags
        );
        uploadSecret(secret)
        reset()
    }

    function handleKeyPress(event: KeyboardEvent) {
        console.log(event.key)
        if (event.key == "Enter") {
            inputTags.push(inputTag)
            inputTag = ""
            inputTags = inputTags
        }
    }

    function reset() {
        inputKey = ""
        inputValue = ""
        inputUrl = ""
        inputTags = []
    }
</script>

<Modal transition={slide} title="New secret" bind:open={openModal} autoclose outsideclose>
    <div class="mb-3">
        <Label>Name</Label>
        <Input bind:value={inputKey} type="text"></Input>
    </div>
    <div class="mb-3">
        <Label>Secret</Label>
        <Input bind:value={inputValue} type="text"></Input>
    </div>
    <div class="mb-3">
        <Label>URL</Label>
        <Input bind:value={inputUrl} type="text"></Input>
    </div>
    <div class="mb-3">
        <Label>Tags</Label>
        <Input bind:value={inputTag} type="text" on:keypress={handleKeyPress}
               placeholder="Type something and press 'Enter' to add a tag"></Input>
        <P slot="left" class="p-2">
            {#each inputTags as tag}
                <Badge class="mr-2" transition={slide} dismissable>
                    #{tag}
                </Badge>
            {/each}
            &nbsp;
        </P>
    </div>
    <svelte:fragment slot="footer">
        <Button color="alternative" data-bs-dismiss="modal">Close</Button>
        <Button color="blue" data-bs-dismiss="modal" on:click={handleSecret}>Save changes</Button>
    </svelte:fragment>
</Modal>