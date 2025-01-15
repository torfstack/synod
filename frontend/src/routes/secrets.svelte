<script lang="ts">
    import {auth} from '$lib/auth'
    import urls from '$lib/config';
    import type {UserCredential} from 'firebase/auth';
    import {Button} from 'flowbite-svelte';
    import KayHeader from "../components/KayHeader.svelte";
    import type {Secret} from "$lib/secret";
    import AddSecretModal from "../components/AddSecretModal.svelte";
    import SearchAddBar from "../components/SearchAddBar.svelte";
    import SecretsList from "../components/SecretsList.svelte";

    interface Props {
        currentUser: UserCredential | null;
    }

    let { currentUser = $bindable() }: Props = $props();
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
        let user = currentUser as UserCredential
        return user.user.getIdToken().then(async token => {
            console.log(token)
            return fetch(urls.backendSecretsUrl, {
                method: "GET",
                credentials: "include",
            })
                .then(resp => resp.json())
                .then(body => {
                    secrets = body as Secret[]
                    console.log("got secrets", secrets)
                })
        });
    }

    async function uploadSecret(s: Secret) {
        let user = currentUser as UserCredential
        return user.user.getIdToken().then(async token => {
            console.log(token)
            await fetch(urls.backendSecretsUrl, {
                method: "POST",
                mode: "cors",
                cache: "no-cache",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include",
                body: JSON.stringify({
                    value: s.value,
                    key: s.key,
                    url: s.url
                })
            })
            await getSecretsFromServer()
        })
    }

    function logout() {
        currentUser = null
        auth.signOut()
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
