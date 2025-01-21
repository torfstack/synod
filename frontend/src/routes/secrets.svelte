<script lang="ts">
    import api from '$lib/api'
    import {Button} from 'flowbite-svelte';
    import KayHeader from "../components/KayHeader.svelte";
    import type {Secret} from "$lib/secret";
    import AddSecretModal from "../components/AddSecretModal.svelte";
    import SearchAddBar from "../components/SearchAddBar.svelte";
    import SecretsList from "../components/SecretsList.svelte";

    interface Props {
        isAuthenticated: boolean;
    }
    let { isAuthenticated = $bindable(false) }: Props = $props();

    let filterValue = $state(""), secrets: Secret[] = $state([]), openModal = $state(false);

    let shown = $derived(secrets.filter((s: Secret) => {
        let trimmed = filterValue.trim()
        const regex = /[A-Z]/
        let hasOnlyLower = trimmed.match(regex) == null
        if (hasOnlyLower) {
            return s.value.toLowerCase().indexOf(trimmed) != -1
        } else {
            return s.value.indexOf(trimmed) != -1
        }
    }))

    async function getSecretsFromServer() {
        api.getSecrets()
          .then(resp => resp.json())
          .then(body => {
              secrets = body as Secret[]
          })
    }

    async function uploadSecret(s: Secret) {
        await api.postSecrets(s)
        await getSecretsFromServer()
    }

    async function logout() {
        await api.deleteAuth()
        isAuthenticated = false
    }

    getSecretsFromServer()
</script>

<KayHeader text="Keeping " />
<div class="p-3">
    <AddSecretModal openModal={openModal} uploadSecret={uploadSecret}/>

    <div class="justify-center flex">
        <Button on:click={logout}>Logout</Button>
    </div>

    <br/>

    <SearchAddBar bind:clickedAdd={openModal} bind:filterValue={filterValue}/>

    <br/>

    <div class="flex justify-center">
        <SecretsList secrets={shown}/>
    </div>
</div>
