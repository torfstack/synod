<script lang="ts">
    import {Secret} from "$lib/secret"
    import {Badge, Button, Modal, P} from "flowbite-svelte";
    import {slide} from "svelte/transition";
    import TextInput from "./input/TextInput.svelte";

    export let uploadSecret: (n: Secret) => Promise<void>;
    export let openModal: boolean;

    let inputKey = "", inputValue = "", inputUrl = "", inputTag = "", inputTags: string[] = []
    let inputKeyError = false, inputValueError = false;

    function handleSecret() {
        if (checkForError()) {
            console.log("some error occured")
            return
        }
        const secret = constructSecret()
        console.log("uploading secret", secret)
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
        if (event.key == "Enter" && inputTag != "") {
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
            inputTags = inputTags.splice(index, 1)
        }
    }
</script>

<Modal transition={slide} title="New secret" bind:open={openModal} outsideclose>
    <div class="mb-3">
        <TextInput label="Name" bind:value={inputKey} bind:error={inputKeyError}
                   errorText="Can not be empty"
                   required={true}/>
    </div>
    <div class="mb-3">
        <TextInput label="Secret" bind:value={inputValue} bind:error={inputValueError}
                   errorText="Can not be empty"
                   required={true}/>
    </div>
    <div class="mb-3">
        <TextInput label="URL" bind:value={inputUrl}/>
    </div>
    <div class="mb-3">
        <TextInput label="Tags" bind:value={inputTag} handleKeyPress={handleKeyPress}
                   placeholder="Type something and press 'Enter' to add a tag"/>
        <P slot="left" class="p-4 h-8">
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