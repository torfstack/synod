<script lang="ts">
    import {Secret} from "$lib/secret"
    import {Badge, Button, Modal, P} from "flowbite-svelte";
    import {slide} from "svelte/transition";
    import TextInput from "./input/TextInput.svelte";

    export let uploadSecret: (n: Secret) => Promise<void>;
    export let openModal: boolean;

    export let inputKey = "", inputValue = "", inputUrl = "", inputTag = "", inputTags: string[] = []
    let inputKeyError = false, inputValueError = false;

    function handleSecret() {
        if (checkForError()) {
            return
        }
        const secret = constructSecret()
        uploadSecret(secret)
        reset()
        openModal = false
    }

    function constructSecret(): Secret {
        return new Secret(
            inputKey,
            inputValue,
            inputUrl,
            inputTags
        );
    }

    function handleKeyPress(event: KeyboardEvent) {
        console.log(event.key)
        if (event.key == "Enter" && inputTag != "" && inputTags.indexOf(inputTag) == -1) {
            inputTags.push(inputTag)
            inputTag = ""
            inputTags = inputTags
        }
    }

    function checkForError(): boolean {
        inputKeyError = inputKey == ""
        inputValueError = inputValue == ""
        return inputKeyError || inputValueError
    }

    function reset() {
        inputKey = ""
        inputValue = ""
        inputUrl = ""
        inputTags = []
    }

    function dismissTag(tag: string) {
        const index = inputTags.indexOf(tag)
        if (index > -1) {
            inputTags.splice(index, 1)
        }
    }
</script>

<Modal bind:open={openModal} outsideclose title="New secret" transition={slide}>
    <div class="mb-3">
        <TextInput bind:error={inputKeyError} bind:value={inputKey} errorText="Can not be empty"
                   label="Name"
                   required={true}/>
    </div>
    <div class="mb-3">
        <TextInput bind:error={inputValueError} bind:value={inputValue} errorText="Can not be empty"
                   label="Secret"
                   required={true}/>
    </div>
    <div class="mb-3">
        <TextInput bind:value={inputUrl} label="URL"/>
    </div>
    <div class="mb-3">
        <TextInput bind:value={inputTag} handleKeyPress={handleKeyPress} label="Tags"
                   placeholder="Type something and press 'Enter' to add a tag"/>
        <P class="p-4 h-8">
            {#each inputTags as tag}
                <Badge class="mr-2 mb-2" transition={slide} dismissable on:dismiss={() => dismissTag(tag)}>
                    #{tag}
                </Badge>
            {/each}
        </P>
    </div>
    <svelte:fragment slot="footer">
        <Button color="alternative" on:click={() => (openModal = false)}>Close</Button>
        <Button color="blue" on:click={handleSecret}>Save changes</Button>
    </svelte:fragment>
</Modal>
