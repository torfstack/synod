<script lang="ts">
    import { Secret } from "$lib/secret"
    import {Button, Input, Label, Modal} from "flowbite-svelte";

    export let uploadSecret: (n: Secret) => Promise<void>;
    export let openModal: boolean;

    let inputKey: string, inputValue: string, inputUrl: string;

    function handleSecret() {
        let secret = new Secret(
            inputKey,
            inputValue,
            inputUrl
        );
        uploadSecret(secret)
    }
</script>

<Modal title="New secret" bind:open={openModal} autoclose outsideclose>
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
    <svelte:fragment slot="footer">
        <Button color="alternative" data-bs-dismiss="modal">Close</Button>
        <Button color="blue" data-bs-dismiss="modal" on:click={handleSecret}>Save changes</Button>
    </svelte:fragment>
</Modal>