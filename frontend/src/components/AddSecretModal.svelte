<script lang="ts">
    import { Secret } from "$lib/secret"
    import {Button, Input, Label, Mark} from "flowbite-svelte";
    let inputKey: string, inputValue: string, inputUrl: string;
    export let uploadSecret: (n: Secret) => Promise<void>;

    function handleSecret() {
        let secret = new Secret(
            inputKey,
            inputValue,
            inputUrl
        );
        uploadSecret(secret)
    }
</script>

<div class="modal fade" id="exampleModal" tabindex="-1" aria-labelledby="exampleModalLabel" >
    <div class="modal-dialog model-dialog-centered">
        <div class="modal-content">
            <div class="modal-header">
                <h1 class="modal-title fs-5" id="exampleModalLabel">New <Mark>secret</Mark></h1>
            </div>
            <div class="modal-body">
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
            </div>
            <div class="modal-footer">
                <Button color="alternative" data-bs-dismiss="modal">Close</Button>
                <Button color="blue" data-bs-dismiss="modal" on:click={handleSecret}>Save changes</Button>
            </div>
        </div>
    </div>
</div>
