<script lang="ts">
    import api from '$lib/api'
    import {Button} from 'flowbite-svelte';
    import KayHeader from "../components/KayHeader.svelte";
    import type {Secret} from "$lib/secret";
    import AddSecretModal from "../components/SecretModal.svelte";
    import SearchAddBar from "../components/SearchAddBar.svelte";
    import SecretsList from "../components/SecretsList.svelte";

    interface Props {
        isAuthenticated: boolean;
    }
    let { isAuthenticated = $bindable(false) }: Props = $props();

    let filterValue = $state(""), secrets: Secret[] = $state([]), openModal = $state(false);
    let selectedSecret = $state<Secret | null>(null)

    let shown = $derived(
        secrets
            .filter((s: Secret) => {
                let trimmed = filterValue.trim()
                const regex = /[A-Z]/
                let hasOnlyLower = trimmed.match(regex) == null
                if (hasOnlyLower) {
                    return s.value.toLowerCase().indexOf(trimmed) != -1
                } else {
                    return s.value.indexOf(trimmed) != -1
                }
            }).sort((a: Secret, b: Secret) => {
                return a.key.localeCompare(b.key)
            })
    )

    async function getSecretsFromServer() {
        await api.getSecrets()
          .then(resp => resp.json())
          .then(body => {
              secrets = body as Secret[]
          })
    }

    async function logout() {
        await api.deleteAuth()
        isAuthenticated = false
    }

    let uploadSecret = $state(async (s: Secret) => {
        await api.postSecrets(s)
        await getSecretsFromServer()
    })

    let selectSecret = $state((s: Secret) => () => {
        selectedSecret = null // how to do this properly? modal is not re-rendered without it
        selectedSecret = s
        openModal = true
    })

    let addNewSecret = $state(() => {
        selectedSecret = null
        openModal = true
    })

    getSecretsFromServer()
</script>

<KayHeader text="Keeping " />
<div class="p-3">
    {#if selectedSecret != null}
        <AddSecretModal
            bind:openModal={openModal}
            inputId={selectedSecret.id}
            inputKey={selectedSecret.key}
            inputTags={selectedSecret.tags}
            inputUrl={selectedSecret.url}
            inputValue={selectedSecret.value}
            uploadSecret={uploadSecret}
        />
    {:else}
        <AddSecretModal
          bind:openModal={openModal}
          uploadSecret={uploadSecret}
        />
    {/if}

    <div class="justify-center flex">
        <Button on:click={logout}>Logout</Button>
    </div>

    <br/>

    <SearchAddBar bind:filterValue={filterValue} clickedAdd={addNewSecret}/>

    <br/>

    <div class="flex justify-center">
        <SecretsList clickedSecret={selectSecret} secrets={shown}/>
    </div>
</div>
